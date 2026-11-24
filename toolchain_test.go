package main

import (
	goversion "go/version"
	"runtime"
	"testing"
)

// minGo is the first Go release that fixes GO-2026-6859 (an HTTP/1 server crash in
// net/http). Scanners read the Go version out of the binary, so a build on an older
// toolchain reports the vulnerability. Move it when go.mod moves.
const minGo = "go1.26.11"

func TestBuiltWithFixedGo(t *testing.T) {
	if goversion.Compare(runtime.Version(), minGo) < 0 {
		t.Fatalf("built with %s, want %s or later (GO-2026-6859)", runtime.Version(), minGo)
	}
}

func TestMinGoRejectsEarlierPatch(t *testing.T) {
	for _, old := range []string{"go1.26.10", "go1.25.9", "go1.26rc1"} {
		if goversion.Compare(old, minGo) >= 0 {
			t.Errorf("%s must be older than %s", old, minGo)
		}
	}
	for _, ok := range []string{"go1.26.11", "go1.26.12", "go1.27.0"} {
		if goversion.Compare(ok, minGo) < 0 {
			t.Errorf("%s must satisfy %s", ok, minGo)
		}
	}
}
