package dashboardserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/dashboardexecute"
	"github.com/turbot/powerpipe/internal/workspace"
	"gopkg.in/olahol/melody.v1"
)

const loadStateServerTestMod = `
mod "local" {
  title = "load state server test"
}

control "c" {
  title = "first"
  sql   = "select 'r' as resource, 'ok' as status, 'fine' as reason"
}
`

var setAppConstants sync.Once

// loadStateMessage is what a client reads of a workspace_load_state message
type loadStateMessage struct {
	Action    string `json:"action"`
	LoadState struct {
		Loading   bool `json:"loading"`
		LoadCount int  `json:"load_count"`
	} `json:"load_state"`
}

// startLoadStateServer loads a workspace, watches it, and serves its dashboard server's WebSocket, as
// powerpipe server does
func startLoadStateServer(t *testing.T) (*workspace.PowerpipeWorkspace, string, string) {
	t.Helper()
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
	dir := t.TempDir()
	modFile := filepath.Join(dir, "mod.pp")
	if err := os.WriteFile(modFile, []byte(loadStateServerTestMod), 0o600); err != nil {
		t.Fatal(err)
	}
	// a disconnecting client cancels its executions
	if dashboardexecute.Executor == nil {
		dashboardexecute.Executor = dashboardexecute.NewDashboardExecutor(nil, nil, backend.SearchPathConfig{})
	}
	ctx := context.Background()
	w, ew := workspace.Load(ctx, dir, workspace.WithLateBinding(false))
	if ew.GetError() != nil {
		t.Fatal(ew.GetError())
	}
	webSocket := melody.New()
	s := &Server{
		mutex:            &sync.Mutex{},
		dashboardClients: map[string]*DashboardClientInfo{},
		webSocket:        webSocket,
		workspace:        w,
	}
	w.RegisterDashboardEventHandler(ctx, s.HandleDashboardEvent)
	if err := w.SetupWatcher(ctx, func(context.Context, error) {}); err != nil {
		t.Fatal(err)
	}
	s.InitAsync(ctx)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		_ = webSocket.HandleRequest(rw, r)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(func() {
		srv.Close()
		_ = webSocket.Close()
		w.Close()
	})
	return w, modFile, "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
}

// connect opens a WebSocket client, which reads its messages into the returned channel
func connect(t *testing.T, url string) chan loadStateMessage {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	messages := make(chan loadStateMessage, 100)
	go func() {
		defer close(messages)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var m loadStateMessage
			if json.Unmarshal(data, &m) == nil && m.Action == "workspace_load_state" {
				messages <- m
			}
		}
	}()
	return messages
}

func nextMessage(t *testing.T, messages chan loadStateMessage, what string, match func(loadStateMessage) bool) loadStateMessage {
	t.Helper()
	timeout := time.After(30 * time.Second)
	for {
		select {
		case m, ok := <-messages:
			if !ok {
				t.Fatalf("the connection closed before %s", what)
			}
			if match(m) {
				return m
			}
		case <-timeout:
			t.Fatalf("no message for %s", what)
		}
	}
}

// a client that connects while the server is idle is told the workspace is not loading, and a reload is then
// pushed to it: loading when it starts, and not loading when it ends
func TestWorkspaceLoadState_PushedToConnectedClient(t *testing.T) {
	w, modFile, url := startLoadStateServer(t)
	// let InitAsync register its handlers
	time.Sleep(200 * time.Millisecond)
	messages := connect(t, url)

	first := nextMessage(t, messages, "the load state on connect", func(loadStateMessage) bool { return true })
	if first.LoadState.Loading || first.LoadState.LoadCount != 1 {
		t.Fatalf("on connect while idle: %+v", first)
	}

	if err := os.WriteFile(modFile, []byte(loadStateServerTestMod+"\n# change\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	nextMessage(t, messages, "the start of the reload", func(m loadStateMessage) bool { return m.LoadState.Loading })
	ended := nextMessage(t, messages, "the end of the reload", func(m loadStateMessage) bool { return !m.LoadState.Loading })
	if ended.LoadState.LoadCount != w.LoadState().LoadCount || ended.LoadState.LoadCount < 2 {
		t.Fatalf("at the end of the reload: %+v, workspace load count %d", ended, w.LoadState().LoadCount)
	}
}

// a client that connects while a reload runs is told at once that the workspace is loading, and told when the
// reload ends, with no request from it
func TestWorkspaceLoadState_ClientConnectingMidReload(t *testing.T) {
	w, modFile, url := startLoadStateServer(t)
	time.Sleep(200 * time.Millisecond)

	// stop the reload inside the load itself, as a slow load of a large workspace would be
	w.LoadLock()
	locked := true
	t.Cleanup(func() {
		if locked {
			w.LoadUnlock()
		}
	})
	if err := os.WriteFile(modFile, []byte(loadStateServerTestMod+"\n# change\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for !w.LoadState().Loading {
		if time.Now().After(deadline) {
			t.Fatal("the reload did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}

	messages := connect(t, url)
	first := nextMessage(t, messages, "the load state on connect", func(loadStateMessage) bool { return true })
	if !first.LoadState.Loading {
		t.Fatalf("a client that connected during a reload was told: %+v", first)
	}

	locked = false
	w.LoadUnlock()
	ended := nextMessage(t, messages, "the end of the reload", func(m loadStateMessage) bool { return !m.LoadState.Loading })
	if ended.LoadState.LoadCount < 2 {
		t.Fatalf("at the end of the reload: %+v", ended)
	}
}

// an idle server sends a connected client nothing more: the load state is pushed on a change, not repeated
func TestWorkspaceLoadState_IdleServerSendsNothingMore(t *testing.T) {
	_, _, url := startLoadStateServer(t)
	time.Sleep(200 * time.Millisecond)
	messages := connect(t, url)
	nextMessage(t, messages, "the load state on connect", func(loadStateMessage) bool { return true })
	select {
	case m := <-messages:
		t.Fatalf("an idle server sent %+v", m)
	case <-time.After(2 * time.Second):
	}
}
