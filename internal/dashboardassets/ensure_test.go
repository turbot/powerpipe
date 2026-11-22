package dashboardassets

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type tarEntry struct {
	name    string
	dir     bool
	content string
}

// buildTarGz writes the entries into an in-memory .tar.gz, in order.
func buildTarGz(t *testing.T, entries []tarEntry) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: 0644, Typeflag: tar.TypeReg, Size: int64(len(e.content))}
		if e.dir {
			hdr.Mode = 0755
			hdr.Typeflag = tar.TypeDir
			hdr.Size = 0
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("write header %q: %v", e.name, err)
		}
		if !e.dir {
			if _, err := tw.Write([]byte(e.content)); err != nil {
				t.Fatalf("write content %q: %v", e.name, err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return &buf
}

// listFiles returns every path under root, relative to it.
func listFiles(t *testing.T, root string) []string {
	t.Helper()
	var paths []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel != "." {
			paths = append(paths, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return paths
}

// The archive the release build makes (tar -C ui/dashboard/build .) unpacks in full,
// including names that merely contain "..".
func TestExtractTarGzUnpacksArchiveInsideAssetsFolder(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "assets")
	if err := os.MkdirAll(dest, 0755); err != nil {
		t.Fatal(err)
	}

	archive := buildTarGz(t, []tarEntry{
		{name: "./", dir: true},
		{name: "./index.html", content: "<html></html>"},
		{name: "./static/", dir: true},
		{name: "./static/js/", dir: true},
		{name: "./static/js/main..chunk.js", content: "js"},
		{name: "./static/../version.json", content: `{"version":"1.0.0"}`},
	})

	if err := extractTarGz(context.Background(), archive, dest); err != nil {
		t.Fatalf("extractTarGz: %v", err)
	}

	for path, want := range map[string]string{
		"index.html":               "<html></html>",
		"static/js/main..chunk.js": "js",
		"version.json":             `{"version":"1.0.0"}`,
	} {
		got, err := os.ReadFile(filepath.Join(dest, path))
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", path, got, want)
		}
	}
	if got := listFiles(t, root); len(got) != 6 {
		t.Errorf("files under the temp dir = %v, want assets/ and its 5 entries", got)
	}
}

// An entry that climbs out of the assets folder, or has an absolute path, is refused with
// an error naming it, and nothing is written outside the folder.
func TestExtractTarGzRefusesEntryOutsideAssetsFolder(t *testing.T) {
	cases := []struct {
		label string
		entry tarEntry
	}{
		{"parent file", tarEntry{name: "../outside.txt", content: "escaped"}},
		{"climb after a folder", tarEntry{name: "./static/../../outside.txt", content: "escaped"}},
		{"parent folder", tarEntry{name: "../outside/", dir: true}},
		{"absolute path", tarEntry{name: "/outside.txt", content: "escaped"}},
		{"empty name", tarEntry{name: "", content: "escaped"}},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			root := t.TempDir()
			dest := filepath.Join(root, "assets")
			if err := os.MkdirAll(dest, 0755); err != nil {
				t.Fatal(err)
			}

			archive := buildTarGz(t, []tarEntry{
				{name: "./", dir: true},
				{name: "./index.html", content: "<html></html>"},
				tc.entry,
			})

			err := extractTarGz(context.Background(), archive, dest)
			if err == nil {
				t.Fatalf("extractTarGz accepted entry %q, want it refused", tc.entry.name)
			}
			if want := `"` + tc.entry.name + `"`; !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not name the entry %s", err.Error(), want)
			}
			// only the folder and the entry before the bad one exist anywhere under the temp dir
			got := listFiles(t, root)
			want := []string{"assets", "assets/index.html"}
			if strings.Join(got, ",") != strings.Join(want, ",") {
				t.Errorf("files under the temp dir = %v, want %v", got, want)
			}
		})
	}
}
