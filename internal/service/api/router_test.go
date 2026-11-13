package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/powerpipe/internal/dashboardserver"
	"gopkg.in/olahol/melody.v1"
)

// newTestRouter returns the router of a server, as `powerpipe server` creates it
func newTestRouter(t *testing.T, listen dashboardserver.ListenType, allowedHosts ...string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	// the router serves the dashboard assets from the install dir
	installDir := app_specific.InstallDir
	app_specific.InstallDir = t.TempDir()
	t.Cleanup(func() { app_specific.InstallDir = installDir })

	api := &APIService{HTTPListen: listen, allowedHosts: allowedHosts, webSocket: melody.New()}
	return api.newRouter()
}

// dial opens a WebSocket to the server as a browser does, and returns the status of the handshake
func dial(t *testing.T, server *httptest.Server, host, origin string) int {
	t.Helper()
	header := http.Header{}
	if host != "" {
		header.Set("Host", host)
	}
	if origin != "" {
		header.Set("Origin", origin)
	}
	conn, resp, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/ws", header)
	if conn != nil {
		_ = conn.Close()
	}
	if resp == nil {
		t.Fatalf("WebSocket handshake failed: %v", err)
	}
	return resp.StatusCode
}

func TestRouter_WebSocket(t *testing.T) {
	local := httptest.NewServer(newTestRouter(t, dashboardserver.ListenTypeLocal, ".pipes.turbot.com", "proxy.example.com"))
	t.Cleanup(local.Close)
	network := httptest.NewServer(newTestRouter(t, dashboardserver.ListenTypeNetwork))
	t.Cleanup(network.Close)
	localNoAllowedHosts := httptest.NewServer(newTestRouter(t, dashboardserver.ListenTypeLocal))
	t.Cleanup(localNoAllowedHosts.Close)
	uiDev := httptest.NewServer(newTestRouter(t, dashboardserver.ListenTypeLocal, "localhost"))
	t.Cleanup(uiDev.Close)

	tests := []struct {
		name         string
		server       *httptest.Server
		host, origin string
		want         int
	}{
		{"not a browser", local, "", "", http.StatusSwitchingProtocols},
		{"dashboard UI", local, "localhost:9033", "http://localhost:9033", http.StatusSwitchingProtocols},
		{"dashboard UI on a loopback IP", local, "127.0.0.1:9033", "http://127.0.0.1:9033", http.StatusSwitchingProtocols},
		{"dashboard UI dev server, localhost allowed", uiDev, "localhost:9033", "http://localhost:3000", http.StatusSwitchingProtocols},
		{"Pipes UI", local, "acme-aws.usea1.dashboard.pipes.turbot.com", "https://pipes.turbot.com", http.StatusSwitchingProtocols},
		{"Pipes UI of a tenant", local, "acme-aws.usea1.dashboard.acme.pipes.turbot.com", "https://acme.pipes.turbot.com", http.StatusSwitchingProtocols},
		{"Pipes API", local, "acme-aws.usea1.dashboard.pipes.turbot.com", "", http.StatusSwitchingProtocols},
		{"proxy that passes the Host", local, "proxy.example.com", "https://proxy.example.com", http.StatusSwitchingProtocols},
		{"proxy that sends 127.0.0.1 as the Host", local, "127.0.0.1:9033", "https://proxy.example.com", http.StatusSwitchingProtocols},
		{"proxy that sends 127.0.0.1 as the Host, not allowed", localNoAllowedHosts, "127.0.0.1:9033", "https://proxy.example.com", http.StatusForbidden},
		{"page on another localhost port", local, "localhost:9033", "http://localhost:8080", http.StatusForbidden},
		{"dashboard UI dev server, localhost not allowed", local, "localhost:9033", "http://localhost:3000", http.StatusForbidden},
		{"page on another site", local, "localhost:9033", "https://evil.example", http.StatusForbidden},
		{"page on another site, network", network, "192.168.1.10:9033", "https://evil.example", http.StatusForbidden},
		{"sandboxed page", local, "localhost:9033", "null", http.StatusForbidden},
		// the page points its own domain at 127.0.0.1, so it is on the same origin as the server
		{"DNS rebinding", local, "evil.example:9033", "http://evil.example:9033", http.StatusForbidden},
		{"host that is not allowed", localNoAllowedHosts, "acme-aws.usea1.dashboard.pipes.turbot.com", "https://pipes.turbot.com", http.StatusForbidden},
		{"any host, network", network, "powerpipe.example.com:9033", "http://powerpipe.example.com:9033", http.StatusSwitchingProtocols},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := dial(t, test.server, test.host, test.origin); got != test.want {
				t.Fatalf("got status %d, want %d", got, test.want)
			}
		})
	}
}

func TestRouter_Host(t *testing.T) {
	local := newTestRouter(t, dashboardserver.ListenTypeLocal, ".pipes.turbot.com", "proxy.example.com")
	network := newTestRouter(t, dashboardserver.ListenTypeNetwork)

	tests := []struct {
		name         string
		router       *gin.Engine
		path, host   string
		want         int
		wantContains string
	}{
		{"localhost", local, "/api/v0/service", "localhost:9033", http.StatusOK, `"UP"`},
		{"loopback IP", local, "/api/v0/service", "127.0.0.1:9033", http.StatusOK, `"UP"`},
		{"IPv6 loopback", local, "/api/v0/service", "[::1]:9033", http.StatusOK, `"UP"`},
		{"no host", local, "/api/v0/service", "", http.StatusOK, `"UP"`},
		{"allowed host", local, "/api/v0/service", "proxy.example.com", http.StatusOK, `"UP"`},
		{"subdomain of an allowed domain", local, "/api/v0/service", "acme-aws.usea1.dashboard.pipes.turbot.com", http.StatusOK, `"UP"`},
		{"host that is not allowed", local, "/api/v0/service", "evil.example:9033", http.StatusForbidden, "--allowed-hosts"},
		{"dashboard assets", local, "/", "evil.example:9033", http.StatusForbidden, "--allowed-hosts"},
		{"unknown path", local, "/no/such/path", "evil.example:9033", http.StatusForbidden, "--allowed-hosts"},
		{"network", network, "/api/v0/service", "evil.example:9033", http.StatusOK, `"UP"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, test.path, nil)
			r.Host = test.host
			w := httptest.NewRecorder()
			test.router.ServeHTTP(w, r)
			if w.Code != test.want || !strings.Contains(w.Body.String(), test.wantContains) {
				t.Fatalf("got status %d and body %q, want status %d and a body that contains %q", w.Code, w.Body.String(), test.want, test.wantContains)
			}
		})
	}
}
