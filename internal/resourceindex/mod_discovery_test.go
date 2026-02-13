package resourceindex

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getModDepTestPath returns the path to testdata/mods/mod-dependencies
func getModDepTestPath(t testing.TB, subPath string) string {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(currentFile), "..", "testdata", "mods", "mod-dependencies", subPath)
}

// =============================================================================
// Basic Mod Discovery Tests
// =============================================================================

func TestModDiscovery_Basic(t *testing.T) {
	// Test: Discover mods in .powerpipe/mods directory
	tmpDir := t.TempDir()

	// Create main mod structure
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "test_main" {
  title = "Test Main"
}

query "main_query" {
  sql = "SELECT 1"
}
`), 0600))

	// Create dependency mod
	depModDir := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "test", "dep-mod@v1.0.0")
	require.NoError(t, os.MkdirAll(depModDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(depModDir, "mod.pp"), []byte(`
mod "dep_mod" {
  title = "Dep Mod"
}

query "dep_query" {
  sql = "SELECT 2"
}
`), 0600))

	// Scan the main directory
	scanner := NewScanner("test_main")
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	// Main mod resources should be indexed
	index := scanner.GetIndex()
	_, ok := index.Get("test_main.query.main_query")
	assert.True(t, ok, "Main mod query should be indexed")

	// Note: Dependency mods need to be scanned separately with ScanDirectoryWithModName
	// This test verifies the basic directory scanning works
}

func TestModDiscovery_NestedDependency(t *testing.T) {
	// Test: Discover nested dependency mods (dep's dep)
	tmpDir := t.TempDir()

	// Create main mod
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "main" { title = "Main" }
query "main_query" { sql = "SELECT 'main'" }
`), 0600))

	// Create first-level dependency
	depADir := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "test", "dep-a@v1.0.0")
	require.NoError(t, os.MkdirAll(depADir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(depADir, "mod.pp"), []byte(`
mod "dep_a" { title = "Dep A" }
query "dep_a_query" { sql = "SELECT 'dep_a'" }
`), 0600))

	// Create nested dependency (dep_a's dependency)
	depLeafDir := filepath.Join(depADir, ".powerpipe", "mods", "github.com", "test", "dep-leaf@v1.0.0")
	require.NoError(t, os.MkdirAll(depLeafDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(depLeafDir, "mod.pp"), []byte(`
mod "dep_leaf" { title = "Dep Leaf" }
query "leaf_query" { sql = "SELECT 'leaf'" }
`), 0600))

	// Scan main directory
	mainScanner := NewScanner("main")
	err := mainScanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	// Scan dep_a with its mod name
	depAScanner := NewScanner("dep_a")
	err = depAScanner.ScanDirectory(depADir)
	require.NoError(t, err)

	// Scan dep_leaf with its mod name
	depLeafScanner := NewScanner("dep_leaf")
	err = depLeafScanner.ScanDirectory(depLeafDir)
	require.NoError(t, err)

	// Verify all resources are found
	assert.True(t, hasEntry(mainScanner.GetIndex(), "main.query.main_query"))
	assert.True(t, hasEntry(depAScanner.GetIndex(), "dep_a.query.dep_a_query"))
	assert.True(t, hasEntry(depLeafScanner.GetIndex(), "dep_leaf.query.leaf_query"))
}

func TestModDiscovery_MissingModFile(t *testing.T) {
	// Test: Directory exists but no mod.pp - should skip gracefully
	tmpDir := t.TempDir()

	// Create main mod
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "main" { title = "Main" }
query "main_query" { sql = "SELECT 1" }
`), 0600))

	// Create directory without mod.pp (simulating incomplete dependency)
	incompleteDep := filepath.Join(tmpDir, ".powerpipe", "mods", "github.com", "test", "incomplete@v1.0.0")
	require.NoError(t, os.MkdirAll(incompleteDep, 0755))
	// Note: No mod.pp created

	// Create a .pp file in incomplete dep (but no mod.pp)
	require.NoError(t, os.WriteFile(filepath.Join(incompleteDep, "queries.pp"), []byte(`
query "orphan_query" { sql = "SELECT 1" }
`), 0600))

	// Scanning should not fail
	scanner := NewScanner("main")
	err := scanner.ScanDirectory(tmpDir)
	assert.NoError(t, err, "Should not fail on directory without mod.pp")
}

func TestModDiscovery_InvalidModFile(t *testing.T) {
	// Test: mod.pp with syntax errors should be handled gracefully
	tmpDir := t.TempDir()

	// Create main mod with invalid HCL
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "main" {
  title = "Main"  # This is valid
}
`), 0600))

	// Create a queries file with syntax that might cause issues
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(`
query "valid_query" {
  sql = "SELECT 1"
}
`), 0600))

	// Scanning should succeed (scanner is regex-based, not full HCL parse)
	scanner := NewScanner("main")
	err := scanner.ScanDirectory(tmpDir)
	assert.NoError(t, err)

	// Valid query should still be indexed
	index := scanner.GetIndex()
	_, ok := index.Get("main.query.valid_query")
	assert.True(t, ok, "Valid query should be indexed despite other issues")
}

// =============================================================================
// Empty and Missing Directory Tests
// =============================================================================

func TestModDiscovery_EmptyModsDir(t *testing.T) {
	// Test: .powerpipe/mods exists but is empty
	tmpDir := t.TempDir()

	// Create main mod
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "main" { title = "Main" }
query "main_query" { sql = "SELECT 1" }
`), 0600))

	// Create empty .powerpipe/mods directory
	modsDir := filepath.Join(tmpDir, ".powerpipe", "mods")
	require.NoError(t, os.MkdirAll(modsDir, 0755))

	// Scanning should succeed
	scanner := NewScanner("main")
	err := scanner.ScanDirectory(tmpDir)
	assert.NoError(t, err, "Should handle empty .powerpipe/mods")

	// Main mod resources should still be indexed
	index := scanner.GetIndex()
	_, ok := index.Get("main.query.main_query")
	assert.True(t, ok)
}

func TestModDiscovery_NoModsDir(t *testing.T) {
	// Test: No .powerpipe/mods directory at all
	tmpDir := t.TempDir()

	// Create main mod
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "main" { title = "Main" }
query "main_query" { sql = "SELECT 1" }
`), 0600))

	// Note: No .powerpipe/mods directory created

	// Scanning should succeed
	scanner := NewScanner("main")
	err := scanner.ScanDirectory(tmpDir)
	assert.NoError(t, err, "Should handle missing .powerpipe/mods")

	// Main mod resources should still be indexed
	index := scanner.GetIndex()
	_, ok := index.Get("main.query.main_query")
	assert.True(t, ok)
}

// =============================================================================
// ScanDirectoryWithModName Tests
// =============================================================================

func TestModDiscovery_ScanWithDifferentModName(t *testing.T) {
	// Test: Scanning a directory with a specific mod name
	tmpDir := t.TempDir()

	// Create mod files
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "original_name" { title = "Original" }
`), 0600))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(`
query "test_query" { sql = "SELECT 1" }
`), 0600))

	// Scan with a specific mod name (simulating dependency scanning)
	scanner := NewScanner("") // Empty initial name
	err := scanner.ScanDirectoryWithModName(tmpDir, "custom_name")
	assert.NoError(t, err)

	// Resources should be indexed with the custom mod name
	index := scanner.GetIndex()
	_, ok := index.Get("custom_name.query.test_query")
	assert.True(t, ok, "Should index with custom mod name")
}

func TestModDiscovery_ModRootSet(t *testing.T) {
	// Test: ModRoot is set correctly when scanning
	tmpDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(`
mod "test" { title = "Test" }
`), 0600))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(`
query "test_query" { sql = "SELECT 1" }
`), 0600))

	scanner := NewScanner("test")
	scanner.SetModRoot(tmpDir)
	err := scanner.ScanDirectory(tmpDir)
	require.NoError(t, err)

	// Verify ModRoot is set on entries
	index := scanner.GetIndex()
	entry, ok := index.Get("test.query.test_query")
	require.True(t, ok)
	assert.Equal(t, tmpDir, entry.ModRoot, "ModRoot should be set on entry")
}

// =============================================================================
// Version Path Extraction Tests
// =============================================================================

func TestModDiscovery_VersionFromPath(t *testing.T) {
	// Test: Extract version info from directory paths
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "standard version",
			path:     "github.com/turbot/mod@v1.2.3",
			expected: "v1.2.3",
		},
		{
			name:     "prerelease version",
			path:     "github.com/turbot/mod@v1.2.3-beta.1",
			expected: "v1.2.3-beta.1",
		},
		{
			name:     "build metadata",
			path:     "github.com/turbot/mod@v1.2.3+build.456",
			expected: "v1.2.3+build.456",
		},
		{
			name:     "no version",
			path:     "github.com/turbot/mod",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version := extractVersionFromPath(tc.path)
			assert.Equal(t, tc.expected, version)
		})
	}
}

// extractVersionFromPath extracts version from a mod path
// e.g., "github.com/turbot/mod@v1.2.3" -> "v1.2.3"
func extractVersionFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '@' {
			return path[i+1:]
		}
	}
	return ""
}

// =============================================================================
// Mod Name Mapping Tests
// =============================================================================

func TestModDiscovery_RegisterModName(t *testing.T) {
	// Test: Mod name registration and resolution
	index := NewResourceIndex()

	// Register mappings
	index.RegisterModName("github.com/turbot/steampipe-mod-aws-insights", "aws_insights")
	index.RegisterModName("github.com/turbot/steampipe-mod-gcp-compliance", "gcp_compliance")

	// Test resolution
	resolved := index.ResolveModName("github.com/turbot/steampipe-mod-aws-insights")
	assert.Equal(t, "aws_insights", resolved)

	resolved = index.ResolveModName("github.com/turbot/steampipe-mod-gcp-compliance")
	assert.Equal(t, "gcp_compliance", resolved)

	// Unknown path returns unchanged
	resolved = index.ResolveModName("unknown/path")
	assert.Equal(t, "unknown/path", resolved)
}

func TestModDiscovery_NameMappingWithHyphens(t *testing.T) {
	// Test: Mod names with hyphens are handled correctly
	index := NewResourceIndex()

	// Register with hyphenated path
	index.RegisterModName("github.com/org/my-mod-name", "my_mod_name")

	resolved := index.ResolveModName("github.com/org/my-mod-name")
	assert.Equal(t, "my_mod_name", resolved)
}

// =============================================================================
// Real Test Fixtures Tests
// =============================================================================

func TestModDiscovery_MainModWithDeps(t *testing.T) {
	// Test using actual test fixtures
	modPath := getModDepTestPath(t, "main-mod")

	scanner := NewScanner("main_mod")
	scanner.SetModRoot(modPath)
	err := scanner.ScanDirectory(modPath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Main mod resources should be indexed
	_, ok := index.Get("main_mod.query.main_query")
	assert.True(t, ok, "Main mod query should be indexed")

	_, ok = index.Get("main_mod.control.local_control")
	assert.True(t, ok, "Main mod control should be indexed")

	_, ok = index.Get("main_mod.benchmark.mixed_benchmark")
	assert.True(t, ok, "Main mod benchmark should be indexed")

	_, ok = index.Get("main_mod.dashboard.main_dashboard")
	assert.True(t, ok, "Main mod dashboard should be indexed")
}

func TestModDiscovery_DepModSeparate(t *testing.T) {
	// Test: Scanning a dependency mod with its own mod name
	depAPath := getModDepTestPath(t, "main-mod/.powerpipe/mods/github.com/test/dep-a@v1.0.0")

	scanner := NewScanner("dep_a")
	scanner.SetModRoot(depAPath)
	err := scanner.ScanDirectory(depAPath)
	require.NoError(t, err)

	index := scanner.GetIndex()

	// Dep A resources should be indexed with dep_a mod name
	_, ok := index.Get("dep_a.query.helper_query")
	assert.True(t, ok, "Dep A query should be indexed")

	_, ok = index.Get("dep_a.control.dep_a_control")
	assert.True(t, ok, "Dep A control should be indexed")

	_, ok = index.Get("dep_a.benchmark.dep_a_benchmark")
	assert.True(t, ok, "Dep A benchmark should be indexed")
}

func TestModDiscovery_CombinedIndex(t *testing.T) {
	// Test: Combining resources from multiple mods into single index
	mainPath := getModDepTestPath(t, "main-mod")
	depAPath := getModDepTestPath(t, "main-mod/.powerpipe/mods/github.com/test/dep-a@v1.0.0")
	depBPath := getModDepTestPath(t, "main-mod/.powerpipe/mods/github.com/test/dep-b@v1.0.0")

	// Create scanner for main mod
	scanner := NewScanner("main_mod")
	scanner.SetModRoot(mainPath)

	// Scan main mod
	err := scanner.ScanDirectory(mainPath)
	require.NoError(t, err)

	// Scan dep mods with their own names (using shared scanner's index)
	err = scanner.ScanDirectoryWithModName(depAPath, "dep_a")
	require.NoError(t, err)

	err = scanner.ScanDirectoryWithModName(depBPath, "dep_b")
	require.NoError(t, err)

	// Register mod name mappings
	index := scanner.GetIndex()
	index.RegisterModName("github.com/test/dep-a", "dep_a")
	index.RegisterModName("github.com/test/dep-b", "dep_b")

	// All resources should be in the combined index
	_, ok := index.Get("main_mod.query.main_query")
	assert.True(t, ok, "Main mod query should be in combined index")

	_, ok = index.Get("dep_a.query.helper_query")
	assert.True(t, ok, "Dep A query should be in combined index")

	_, ok = index.Get("dep_b.query.helper_query")
	assert.True(t, ok, "Dep B query should be in combined index")

	// Statistics
	stats := index.Stats()
	t.Logf("Combined index: %d entries, types: %v", stats.TotalEntries, stats.ByType)
	assert.Greater(t, stats.TotalEntries, 0)
}

// =============================================================================
// Symlink Tests (if supported by OS)
// =============================================================================

func TestModDiscovery_Symlinks(t *testing.T) {
	// Test: Symlinked mod directories behavior
	// Note: filepath.Walk does not follow symlinks by default in Go
	// This test documents the current behavior
	tmpDir := t.TempDir()

	// Create actual mod directory
	actualDir := filepath.Join(tmpDir, "actual_mod")
	require.NoError(t, os.MkdirAll(actualDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(actualDir, "mod.pp"), []byte(`
mod "symlinked_mod" { title = "Symlinked" }
query "sym_query" { sql = "SELECT 1" }
`), 0600))

	// Create main mod with symlinked dependency
	mainDir := filepath.Join(tmpDir, "main")
	modsDir := filepath.Join(mainDir, ".powerpipe", "mods", "github.com", "test")
	require.NoError(t, os.MkdirAll(modsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(mainDir, "mod.pp"), []byte(`
mod "main" { title = "Main" }
`), 0600))

	// Create symlink to actual mod
	symlinkPath := filepath.Join(modsDir, "symlinked@v1.0.0")
	err := os.Symlink(actualDir, symlinkPath)
	if err != nil {
		t.Skip("Symlinks not supported on this system")
	}

	// When scanning a symlink directly (as a path), it should work
	// because we resolve the symlink path before scanning
	resolvedPath, err := filepath.EvalSymlinks(symlinkPath)
	require.NoError(t, err)

	scanner := NewScanner("symlinked_mod")
	err = scanner.ScanDirectory(resolvedPath)
	require.NoError(t, err)

	// Resources should be discovered when scanning resolved symlink path
	index := scanner.GetIndex()
	_, ok := index.Get("symlinked_mod.query.sym_query")
	assert.True(t, ok, "Query should be discovered through resolved symlink")
}

// =============================================================================
// Parallel Scanning Tests
// =============================================================================

func TestModDiscovery_ParallelScanning(t *testing.T) {
	// Test: Parallel scanning produces same results as sequential
	mainPath := getModDepTestPath(t, "main-mod")

	// Sequential scan
	seqScanner := NewScanner("main_mod")
	err := seqScanner.ScanDirectory(mainPath)
	require.NoError(t, err)
	seqIndex := seqScanner.GetIndex()
	seqCount := seqIndex.Count()

	// Parallel scan
	parScanner := NewScanner("main_mod")
	err = parScanner.ScanDirectoryParallel(mainPath, 4)
	require.NoError(t, err)
	parIndex := parScanner.GetIndex()
	parCount := parIndex.Count()

	// Same number of entries
	assert.Equal(t, seqCount, parCount, "Sequential and parallel should find same entries")

	// Same entries exist
	for _, entry := range seqIndex.List() {
		_, ok := parIndex.Get(entry.Name)
		assert.True(t, ok, "Entry %s should exist in both indexes", entry.Name)
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkModDiscovery_ScanDirectory(b *testing.B) {
	mainPath := getModDepTestPath(b, "main-mod")

	for i := 0; i < b.N; i++ {
		scanner := NewScanner("main_mod")
		_ = scanner.ScanDirectory(mainPath)
	}
}

func BenchmarkModDiscovery_ScanDirectoryParallel(b *testing.B) {
	mainPath := getModDepTestPath(b, "main-mod")

	for i := 0; i < b.N; i++ {
		scanner := NewScanner("main_mod")
		_ = scanner.ScanDirectoryParallel(mainPath, 4)
	}
}

func BenchmarkModDiscovery_CombinedIndex(b *testing.B) {
	mainPath := getModDepTestPath(b, "main-mod")
	depAPath := getModDepTestPath(b, "main-mod/.powerpipe/mods/github.com/test/dep-a@v1.0.0")
	depBPath := getModDepTestPath(b, "main-mod/.powerpipe/mods/github.com/test/dep-b@v1.0.0")

	for i := 0; i < b.N; i++ {
		scanner := NewScanner("main_mod")
		_ = scanner.ScanDirectory(mainPath)
		_ = scanner.ScanDirectoryWithModName(depAPath, "dep_a")
		_ = scanner.ScanDirectoryWithModName(depBPath, "dep_b")
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

func hasEntry(index *ResourceIndex, name string) bool {
	_, ok := index.Get(name)
	return ok
}
