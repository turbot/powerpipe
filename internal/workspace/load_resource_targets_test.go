package workspace_test

import (
	"context"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/powerpipe/internal/cmdconfig"
	"github.com/turbot/powerpipe/internal/workspace"
)

var setAppConstants sync.Once

const targetsTestMod = `
mod "local" {
  title = "resource targets test"
}

locals {
  tags = { service = "AWS/S3" }
}

benchmark "top" {
  title    = "Top"
  children = [benchmark.child, control.c1]
  tags     = local.tags
}

benchmark "child" {
  title    = "Child, after ${control.c1.title}"
  children = [control.c2]
}

control "c1" {
  title = "C1"
  query = query.q1
}

control "c2" {
  title         = "C2"
  sql           = "select 'r' as resource, 'ok' as status, 'reason' as reason"
  documentation = file("./docs/c2.md")
}

control "unused" {
  sql = "select 3"
}

// an unused benchmark that shares control c1 with benchmark top
benchmark "other" {
  children = [control.c1, control.unused]
}

query "q1" {
  sql = "select 'r' as resource, 'ok' as status, 'reason' as reason"
}

dashboard "d" {
  card {
    sql = "select 1"
  }
}
`

// writeTestMod writes a mod with the given extra source and returns its path
func writeTestMod(t *testing.T, extra string) string {
	t.Helper()
	setAppConstants.Do(func() {
		viper.Set("main.version", "0.0.0")
		cmdconfig.SetAppSpecificConstants()
	})
	dir := t.TempDir()
	files := map[string]string{
		"mod.pp":     targetsTestMod,
		"extra.pp":   extra,
		"docs/c2.md": "# C2 docs",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func loadTestMod(t *testing.T, dir string, targets ...parse.ResourceTarget) (*workspace.PowerpipeWorkspace, error) {
	t.Helper()
	w, ew := workspace.Load(context.Background(), dir, workspace.WithResourceTargets(targets), workspace.WithLateBinding(false))
	return w, ew.GetError()
}

// resourceJson returns the JSON of every resource except the mod, keyed by name, with two changes:
//   - the references are sorted: their order comes from a map, so it changes between loads
//   - the path is removed, and returned separately: in a partial load, a resource shared with a benchmark
//     that is not loaded does not have the path through that benchmark (a benchmark run does not use paths)
func resourceJson(t *testing.T, w *workspace.PowerpipeWorkspace) (resources map[string]string, paths map[string]string) {
	t.Helper()
	resources = make(map[string]string)
	paths = make(map[string]string)
	err := w.GetModResources().WalkResources(func(r modconfig.HclResource) (bool, error) {
		if _, isMod := r.(*modconfig.Mod); isMod {
			return true, nil
		}
		b, err := json.Marshal(r)
		if err != nil {
			return false, err
		}
		var fields map[string]any
		if err := json.Unmarshal(b, &fields); err != nil {
			return false, err
		}
		if refs, ok := fields["references"].([]any); ok {
			slices.SortFunc(refs, func(a, b any) int {
				aj, _ := json.Marshal(a)
				bj, _ := json.Marshal(b)
				return strings.Compare(string(aj), string(bj))
			})
		}
		path, _ := json.Marshal(fields["path"])
		delete(fields, "path")
		b, err = json.Marshal(fields)
		if err != nil {
			return false, err
		}
		resources[r.Name()] = string(b)
		paths[r.Name()] = string(path)
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return resources, paths
}

func TestLoad_ResourceTargets(t *testing.T) {
	dir := writeTestMod(t, "")
	full, err := loadTestMod(t, dir)
	if err != nil {
		t.Fatal(err)
	}
	partial, err := loadTestMod(t, dir, parse.ResourceTarget{Name: "benchmark.top"})
	if err != nil {
		t.Fatal(err)
	}

	fullResources, fullPaths := resourceJson(t, full)
	partialResources, partialPaths := resourceJson(t, partial)
	want := []string{"local.benchmark.child", "local.benchmark.top", "local.control.c1", "local.control.c2", "local.local.tags", "local.query.q1"}
	if got := slices.Sorted(maps.Keys(partialResources)); !slices.Equal(got, want) {
		t.Fatalf("partial load resources\n got: %v\nwant: %v", got, want)
	}
	for name, partialJson := range partialResources {
		if partialJson != fullResources[name] {
			t.Errorf("%s differs from the full load\npartial: %s\n   full: %s", name, partialJson, fullResources[name])
		}
	}
	if _, ok := fullResources["local.control.unused"]; !ok {
		t.Fatalf("expected the full load to have every resource")
	}

	// control c1 is shared with benchmark other, which is not loaded: only its path through benchmark top is known
	wantPath := `[["mod.local","local.benchmark.top","local.control.c1"]]`
	if partialPaths["local.control.c1"] != wantPath || !strings.Contains(fullPaths["local.control.c1"], "local.benchmark.other") {
		t.Fatalf("unexpected paths of control c1\npartial: %s\n   full: %s", partialPaths["local.control.c1"], fullPaths["local.control.c1"])
	}
	for name, path := range partialPaths {
		if name != "local.control.c1" && path != fullPaths[name] {
			t.Errorf("%s path differs\npartial: %s\n   full: %s", name, path, fullPaths[name])
		}
	}
}

func TestLoad_ResourceTargets_ModQualified(t *testing.T) {
	dir := writeTestMod(t, "")
	w, err := loadTestMod(t, dir, parse.ResourceTarget{Mod: "local", Name: "control.c1"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"local.control.c1", "local.local.tags", "local.query.q1"}
	resources, _ := resourceJson(t, w)
	if got := slices.Sorted(maps.Keys(resources)); !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// if a target is not defined, every resource is loaded, so that the error for the missing target is the same
func TestLoad_ResourceTargets_MissingTarget(t *testing.T) {
	dir := writeTestMod(t, "")
	w, err := loadTestMod(t, dir, parse.ResourceTarget{Mod: "local", Name: "benchmark.missing"})
	if err != nil {
		t.Fatal(err)
	}
	if resources, _ := resourceJson(t, w); resources["local.control.unused"] == "" {
		t.Fatal("expected every resource to be loaded")
	}
}

// a broken resource that the target does not use does not stop the load. A syntax error anywhere still does.
func TestLoad_ResourceTargets_UnusedErrors(t *testing.T) {
	dir := writeTestMod(t, `control "broken" { query = query.does_not_exist }`)
	if _, err := loadTestMod(t, dir); err == nil {
		t.Fatal("expected the full load to fail")
	}
	if _, err := loadTestMod(t, dir, parse.ResourceTarget{Name: "benchmark.top"}); err != nil {
		t.Fatalf("expected the partial load to succeed: %v", err)
	}

	syntaxError := writeTestMod(t, `control "broken" {`)
	if _, err := loadTestMod(t, syntaxError, parse.ResourceTarget{Name: "benchmark.top"}); err == nil {
		t.Fatal("expected a syntax error to fail the partial load")
	}
}
