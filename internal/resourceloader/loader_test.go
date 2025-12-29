package resourceloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"

	pparse "github.com/turbot/powerpipe/internal/parse"
)

func init() {
	// Set up app-specific constants required for mod loading
	app_specific.AppName = "powerpipe"
	app_specific.ModDataExtensions = []string{".pp", ".sp"}
	app_specific.VariablesExtensions = []string{".ppvars", ".spvars"}
	app_specific.AutoVariablesExtensions = []string{".auto.ppvars", ".auto.spvars"}
	app_specific.DefaultVarsFileName = "powerpipe.ppvars"
	app_specific.LegacyDefaultVarsFileName = "steampipe.spvars"
	app_specific.WorkspaceIgnoreFile = ".powerpipeignore"
	app_specific.WorkspaceDataDir = ".powerpipe"

	// Set up app-specific functions required for mod loading
	modconfig.AppSpecificNewModResourcesFunc = resources.NewModResources
	parse.ModDecoderFunc = pparse.NewPowerpipeModDecoder
	parse.AppSpecificGetResourceSchemaFunc = pparse.GetResourceSchema
}

func TestLoader_LoadQuery(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()
	resource, err := loader.Load(ctx, "testmod.query.simple")
	require.NoError(t, err)
	assert.NotNil(t, resource)

	query, ok := resource.(*resources.Query)
	require.True(t, ok)
	assert.Equal(t, "simple", query.ShortName)
	assert.NotNil(t, query.SQL)
	assert.Equal(t, "SELECT 1", *query.SQL)
}

func TestLoader_LoadControl(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()
	resource, err := loader.Load(ctx, "testmod.control.check_one")
	require.NoError(t, err)
	assert.NotNil(t, resource)

	control, ok := resource.(*resources.Control)
	require.True(t, ok)
	assert.Equal(t, "check_one", control.ShortName)
	assert.NotNil(t, control.SQL)
}

func TestLoader_CacheHit(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()

	// First load - miss
	_, err := loader.Load(ctx, "testmod.query.simple")
	require.NoError(t, err)
	stats := loader.cache.Stats()
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, int64(0), stats.Hits)

	// Second load - hit
	_, err = loader.Load(ctx, "testmod.query.simple")
	require.NoError(t, err)
	stats = loader.cache.Stats()
	assert.Equal(t, int64(1), stats.Hits)
}

func TestLoader_NotFound(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()
	_, err := loader.Load(ctx, "testmod.query.nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLoader_Preload(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	names := []string{"testmod.query.q1", "testmod.query.q2", "testmod.query.q3"}
	err := loader.Preload(context.Background(), names)
	require.NoError(t, err)

	// All should now be in cache
	for _, name := range names {
		_, ok := loader.cache.GetResource(name)
		assert.True(t, ok, "expected %s to be in cache", name)
	}
}

func TestLoader_Stats(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()

	// Load some resources
	_, _ = loader.Load(ctx, "testmod.query.simple")
	_, _ = loader.Load(ctx, "testmod.query.q1")

	stats := loader.Stats()
	assert.Equal(t, int64(2), stats.LoadCount)
	assert.Greater(t, stats.AvgParseTime.Nanoseconds(), int64(0))
}

func TestLoader_Invalidate(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()

	// Load a resource
	_, err := loader.Load(ctx, "testmod.query.simple")
	require.NoError(t, err)

	// Verify it's in cache
	_, ok := loader.cache.GetResource("testmod.query.simple")
	assert.True(t, ok)

	// Invalidate
	loader.Invalidate("testmod.query.simple")

	// Verify it's not in cache
	_, ok = loader.cache.GetResource("testmod.query.simple")
	assert.False(t, ok)
}

func TestLoader_Clear(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	ctx := context.Background()

	// Load some resources
	_, _ = loader.Load(ctx, "testmod.query.simple")
	_, _ = loader.Load(ctx, "testmod.query.q1")

	assert.Equal(t, 2, loader.cache.Len())

	// Clear
	loader.Clear()

	assert.Equal(t, 0, loader.cache.Len())
	assert.Equal(t, int64(0), loader.Stats().LoadCount)
}

func TestLoader_PreloadByType(t *testing.T) {
	modPath, mod := setupTestMod(t)
	loader := setupTestLoader(t, modPath, mod)

	err := loader.PreloadByType(context.Background(), "query")
	require.NoError(t, err)

	// All queries should be in cache
	entries := loader.index.GetByType("query")
	for _, entry := range entries {
		_, ok := loader.cache.GetResource(entry.Name)
		assert.True(t, ok, "expected %s to be in cache", entry.Name)
	}
}

// Helper functions

func setupTestMod(t testing.TB) (string, *modconfig.Mod) {
	tmpDir := t.TempDir()

	// Create mod.pp
	modContent := `mod "testmod" {
  title = "Test Mod"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "mod.pp"), []byte(modContent), 0644))

	// Create queries.pp
	queriesContent := `query "simple" {
  sql = "SELECT 1"
}

query "q1" {
  sql = "SELECT 1"
}

query "q2" {
  sql = "SELECT 2"
}

query "q3" {
  sql = "SELECT 3"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "queries.pp"), []byte(queriesContent), 0644))

	// Create controls.pp
	controlsContent := `control "check_one" {
  sql = "SELECT 'ok' as status, 'test' as reason"
  title = "Check One"
}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "controls.pp"), []byte(controlsContent), 0644))

	// Create a minimal mod using the constructor
	mod := modconfig.NewMod("testmod", tmpDir, hcl.Range{})

	return tmpDir, mod
}

func setupTestLoader(t testing.TB, modPath string, mod *modconfig.Mod) *Loader {
	// Create index with test resources
	index := resourceindex.NewResourceIndex()
	index.ModName = "testmod"

	queriesFile := filepath.Join(modPath, "queries.pp")
	controlsFile := filepath.Join(modPath, "controls.pp")

	// Add query entries - calculate approximate byte offsets
	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.simple",
		ShortName: "simple",
		FileName:  queriesFile,
		StartLine: 1,
		EndLine:   3,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q1",
		ShortName: "q1",
		FileName:  queriesFile,
		StartLine: 5,
		EndLine:   7,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q2",
		ShortName: "q2",
		FileName:  queriesFile,
		StartLine: 9,
		EndLine:   11,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "query",
		Name:      "testmod.query.q3",
		ShortName: "q3",
		FileName:  queriesFile,
		StartLine: 13,
		EndLine:   15,
		ModName:   "testmod",
	})

	index.Add(&resourceindex.IndexEntry{
		Type:      "control",
		Name:      "testmod.control.check_one",
		ShortName: "check_one",
		FileName:  controlsFile,
		StartLine: 1,
		EndLine:   4,
		ModName:   "testmod",
	})

	// Create cache
	cache := resourcecache.NewResourceCache(resourcecache.DefaultConfig())

	// Create loader
	return NewLoader(index, cache, mod, modPath)
}
