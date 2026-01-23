// Package workspace provides lazy loading workspace capabilities.
package workspace

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/dashboardevents"
	"github.com/turbot/powerpipe/internal/resourcecache"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resourceloader"
	"github.com/turbot/powerpipe/internal/resources"
)

// LazyWorkspace wraps PowerpipeWorkspace with lazy loading capabilities.
// Instead of loading all resources at startup, it builds an index of resource
// metadata and loads resources on-demand when accessed.
//
// Hybrid Mode: Uses lazy loading for browsing (fast startup), but falls back
// to eager loading when execution is needed (proper reference resolution).
type LazyWorkspace struct {
	*PowerpipeWorkspace

	// Index of all resources (always loaded at startup)
	index *resourceindex.ResourceIndex

	// Cache for parsed resources
	cache *resourcecache.ResourceCache

	// Loader for on-demand parsing
	loader *resourceloader.Loader

	// Resolver for dependencies
	resolver *resourceloader.DependencyResolver

	// Config
	config LazyLoadConfig

	// Lazy mod resources accessor
	lazyResources *LazyModResources

	// Eager workspace for execution (loaded on-demand when first execution is requested)
	eagerWorkspace *PowerpipeWorkspace
	eagerLoadOnce  sync.Once
	eagerLoadErr   error

	// Path for eager loading
	workspacePath string

	// Per-benchmark resolution tracking to prevent concurrent modifications
	// Key: benchmark name, Value: *benchmarkResolution
	resolvedBenchmarks sync.Map

	// Background resolution
	backgroundResolver *BackgroundResolver
	updateListeners    []UpdateListener
	updateListenersMu  sync.RWMutex
	fullyResolved      bool
}

// benchmarkResolution tracks the resolution state of a benchmark
type benchmarkResolution struct {
	once sync.Once
	err  error
}

// LazyLoadConfig configures lazy loading behavior.
type LazyLoadConfig struct {
	// Maximum cache memory in bytes (default: 50MB)
	MaxCacheMemory int64

	// Whether to preload frequently accessed resources
	EnablePreload bool

	// Resources to preload (e.g., top-level benchmarks)
	PreloadPatterns []string
}

// DefaultLazyLoadConfig returns default configuration.
func DefaultLazyLoadConfig() LazyLoadConfig {
	return LazyLoadConfig{
		MaxCacheMemory:  50 * 1024 * 1024, // 50MB
		EnablePreload:   false,
		PreloadPatterns: []string{},
	}
}

// NewLazyWorkspace creates a lazy-loading workspace.
// This is much faster than the standard Load() as it only builds an index
// of resource metadata without fully parsing any resources.
func NewLazyWorkspace(ctx context.Context, workspacePath string, config LazyLoadConfig) (*LazyWorkspace, error) {
	// Scan mod.pp to get mod name and info
	modName, modFullName, modTitle, err := scanModInfo(workspacePath)
	if err != nil {
		return nil, fmt.Errorf("scanning mod info: %w", err)
	}

	// Build index from files (fast scan, no full parse)
	index, err := buildResourceIndex(ctx, workspacePath, modName)
	if err != nil {
		return nil, fmt.Errorf("building index: %w", err)
	}

	// Set mod info on index
	index.ModName = modName
	index.ModFullName = modFullName
	index.ModTitle = modTitle

	// Create cache with memory limit
	cacheConfig := resourcecache.CacheConfig{
		MaxMemoryBytes: config.MaxCacheMemory,
	}
	cache := resourcecache.NewResourceCache(cacheConfig)

	// Create minimal mod for the loader
	mod := modconfig.NewMod(modName, workspacePath, hcl.Range{})

	// Create loader
	loader := resourceloader.NewLoader(index, cache, mod, workspacePath)

	// Create resolver
	resolver := resourceloader.NewDependencyResolver(index, loader)

	// Create base PowerpipeWorkspace with minimal initialization
	pw := NewPowerpipeWorkspace(workspacePath)
	pw.Mod = mod

	lw := &LazyWorkspace{
		PowerpipeWorkspace: pw,
		index:              index,
		cache:              cache,
		loader:             loader,
		resolver:           resolver,
		config:             config,
		workspacePath:      workspacePath,
	}

	// Create lazy mod resources accessor
	lw.lazyResources = NewLazyModResources(lw)

	// Set the resource provider on the loader for reference resolution
	loader.SetResourceProvider(lw)

	// Optional preload
	if config.EnablePreload && len(config.PreloadPatterns) > 0 {
		lw.preloadResources(ctx, config.PreloadPatterns)
	}

	return lw, nil
}

// GetWorkspaceForExecution returns a fully-loaded workspace for execution.
// This uses eager loading with proper reference resolution, which is needed
// for controls that reference queries. The eager workspace is loaded once
// on first request and cached for subsequent executions.
func (lw *LazyWorkspace) GetWorkspaceForExecution(ctx context.Context) (*PowerpipeWorkspace, error) {
	lw.eagerLoadOnce.Do(func() {
		// Load the workspace eagerly using the standard Load function
		// This does full HCL parsing with proper reference resolution
		ew, errAndWarnings := Load(ctx, lw.workspacePath)
		if errAndWarnings.GetError() != nil {
			lw.eagerLoadErr = errAndWarnings.GetError()
			return
		}

		// Copy event handlers from the lazy workspace to the eager workspace
		// This ensures dashboard events from execution are properly routed to the server
		for _, handler := range lw.PowerpipeWorkspace.dashboardEventHandlers {
			ew.RegisterDashboardEventHandler(ctx, handler)
		}

		lw.eagerWorkspace = ew
	})

	if lw.eagerLoadErr != nil {
		return nil, lw.eagerLoadErr
	}
	return lw.eagerWorkspace, nil
}

// buildResourceIndex scans the workspace and builds a resource index.
func buildResourceIndex(ctx context.Context, workspacePath, modName string) (*resourceindex.ResourceIndex, error) {
	scanner := resourceindex.NewScanner(modName)
	// Set the mod root for the main workspace (needed for file() function resolution)
	scanner.SetModRoot(workspacePath)

	// Scan the main workspace directory
	if err := scanner.ScanDirectoryParallel(workspacePath, 0); err != nil {
		return nil, err
	}

	// Scan dependency mods in .powerpipe/mods directory
	// Each dependency mod needs to be scanned with its own mod name
	modsDir := filepath.Join(workspacePath, ".powerpipe", "mods")
	if info, err := os.Stat(modsDir); err == nil && info.IsDir() {
		if err := scanDependencyMods(scanner, modsDir); err != nil {
			return nil, fmt.Errorf("scanning dependency mods: %w", err)
		}
	}

	// Mark top-level resources and set parent names
	scanner.MarkTopLevelResources()
	scanner.SetParentNames()
	scanner.ComputePaths()

	return scanner.GetIndex(), nil
}

// ErrDuplicateModVersions is returned when multiple versions of the same mod are detected
var ErrDuplicateModVersions = errors.New("duplicate mod versions detected")

// scanDependencyMods walks the mods directory and scans each mod with its own mod name.
func scanDependencyMods(scanner *resourceindex.Scanner, modsDir string) error {
	// Track seen mod names to detect duplicates (diamond dependency issues)
	seenModNames := make(map[string]string) // mod name -> first seen path

	// Walk through the mods directory looking for mod.pp or mod.sp files
	return filepath.Walk(modsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for mod.pp or mod.sp files (both are valid mod definition files)
		if info.IsDir() || (info.Name() != "mod.pp" && info.Name() != "mod.sp") {
			return nil
		}

		// Found a mod definition file - extract mod name and scan its directory
		modDir := filepath.Dir(path)
		depModName, _, _, err := scanModInfo(modDir)
		if err != nil {
			return fmt.Errorf("scanning mod info from %s: %w", path, err)
		}

		// Check for duplicate mod names (different versions of same mod)
		// This indicates a diamond dependency that lazy loading can't handle correctly
		if firstPath, exists := seenModNames[depModName]; exists {
			slog.Info("Duplicate mod versions detected - falling back to eager loading",
				"mod", depModName,
				"path1", firstPath,
				"path2", modDir)
			return ErrDuplicateModVersions
		}
		seenModNames[depModName] = modDir

		// Extract the full mod path from the directory structure
		// e.g., ".powerpipe/mods/github.com/turbot/steampipe-mod-aws-insights@v1.2.0"
		// -> "github.com/turbot/steampipe-mod-aws-insights"
		relPath, _ := filepath.Rel(modsDir, modDir)
		fullModPath := strings.Split(relPath, "@")[0] // Remove version suffix

		// Register the mapping from full path to short name
		scanner.GetIndex().RegisterModName(fullModPath, depModName)

		// Scan this mod directory with the correct mod name
		if err := scanner.ScanDirectoryWithModName(modDir, depModName); err != nil {
			return fmt.Errorf("scanning mod %s: %w", depModName, err)
		}

		// Skip subdirectories of this mod (we've already scanned it)
		return filepath.SkipDir
	})
}

// scanModInfo extracts mod name, full name, and title from mod.pp or mod.sp.
func scanModInfo(workspacePath string) (modName, modFullName, modTitle string, err error) {
	// Try mod.pp first, then mod.sp
	modFilePath := filepath.Join(workspacePath, "mod.pp")
	file, err := os.Open(modFilePath)
	if err != nil && os.IsNotExist(err) {
		// Try mod.sp as fallback
		modFilePath = filepath.Join(workspacePath, "mod.sp")
		file, err = os.Open(modFilePath)
	}

	if err != nil {
		if os.IsNotExist(err) {
			// Default to directory name if no mod.pp or mod.sp
			modName = filepath.Base(workspacePath)
			modFullName = "mod." + modName
			return modName, modFullName, modName, nil
		}
		return "", "", "", err
	}
	defer file.Close()

	// Simple regex-based extraction
	modBlockRegex := regexp.MustCompile(`^\s*mod\s+"([^"]+)"`)
	titleRegex := regexp.MustCompile(`^\s*title\s*=\s*"([^"]*)"`)

	scanner := bufio.NewScanner(file)
	inModBlock := false
	braceDepth := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Only match the first mod block (not nested mod blocks in require section)
		if !inModBlock {
			if matches := modBlockRegex.FindStringSubmatch(line); len(matches) >= 2 {
				modName = matches[1]
				modFullName = "mod." + modName
				inModBlock = true
				// Count braces on this line
				braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
				continue
			}
		}

		if inModBlock {
			// Track brace depth to know when we exit the mod block
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")

			if matches := titleRegex.FindStringSubmatch(line); len(matches) >= 2 {
				modTitle = matches[1]
			}
			// Stop when we've exited the top-level mod block
			if braceDepth <= 0 {
				break
			}
		}
	}

	if modName == "" {
		modName = filepath.Base(workspacePath)
		modFullName = "mod." + modName
	}
	if modTitle == "" {
		modTitle = modName
	}

	return modName, modFullName, modTitle, scanner.Err()
}

// preloadResources loads resources matching patterns in the background.
func (lw *LazyWorkspace) preloadResources(ctx context.Context, patterns []string) {
	// Find matching resources
	var names []string
	for _, pattern := range patterns {
		matches := lw.findByPattern(pattern)
		names = append(names, matches...)
	}

	if len(names) == 0 {
		return
	}

	// Preload with dependencies in background
	go func() {
		_ = lw.loader.PreloadWithDependencies(ctx, names, resourceloader.PreloadOptions{
			IncludeDependencies: true,
			MaxConcurrency:      10,
		})
	}()
}

// findByPattern finds resource names matching a pattern.
func (lw *LazyWorkspace) findByPattern(pattern string) []string {
	var matches []string

	for _, entry := range lw.index.List() {
		// Simple pattern matching - support * wildcard
		if matchPattern(pattern, entry.Name) || matchPattern(pattern, entry.ShortName) {
			matches = append(matches, entry.Name)
		}
	}

	return matches
}

// matchPattern does simple wildcard matching.
func matchPattern(pattern, name string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(name, pattern[1:]) {
		return true
	}
	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(name, pattern[:len(pattern)-1]) {
		return true
	}
	return pattern == name
}

// GetResource retrieves a resource by parsed name, loading it on-demand if needed.
// This implements the modconfig.ResourceProvider interface.
// For benchmarks, this resolves children so they're available for execution.
func (lw *LazyWorkspace) GetResource(parsedName *modconfig.ParsedResourceName) (modconfig.HclResource, bool) {
	ctx := context.Background()

	// Build full name
	modName := parsedName.Mod
	if modName == "" {
		modName = lw.index.ModName
	}

	// Convert full mod path (e.g., "github.com/turbot/steampipe-mod-aws-insights")
	// to short name (e.g., "aws_insights") for index lookup
	modName = lw.index.ResolveModName(modName)

	fullName := fmt.Sprintf("%s.%s.%s", modName, parsedName.ItemType, parsedName.Name)

	// For benchmarks, use LoadBenchmarkForExecution to ensure children are resolved
	if parsedName.ItemType == "benchmark" {
		resource, err := lw.LoadBenchmarkForExecution(ctx, fullName)
		if err != nil {
			return nil, false
		}
		return resource.(modconfig.HclResource), true
	}

	// Try to load from cache or disk
	resource, err := lw.loader.Load(ctx, fullName)
	if err != nil {
		return nil, false
	}

	return resource, true
}

// GetAvailableDashboardsFromIndex builds the available dashboards payload
// without loading any resources - uses only the index.
func (lw *LazyWorkspace) GetAvailableDashboardsFromIndex() *resourceindex.AvailableDashboardsPayload {
	return lw.index.BuildAvailableDashboardsPayload()
}

// GetIndex returns the resource index.
func (lw *LazyWorkspace) GetIndex() *resourceindex.ResourceIndex {
	return lw.index
}

// GetLoader returns the resource loader.
func (lw *LazyWorkspace) GetLoader() *resourceloader.Loader {
	return lw.loader
}

// GetCache returns the resource cache.
func (lw *LazyWorkspace) GetCache() *resourcecache.ResourceCache {
	return lw.cache
}

// GetResolver returns the dependency resolver.
func (lw *LazyWorkspace) GetResolver() *resourceloader.DependencyResolver {
	return lw.resolver
}

// GetLazyModResources returns the lazy mod resources accessor.
func (lw *LazyWorkspace) GetLazyModResources() *LazyModResources {
	return lw.lazyResources
}

// LoadDashboard loads a dashboard and all its children on-demand.
func (lw *LazyWorkspace) LoadDashboard(ctx context.Context, name string) (*resources.Dashboard, error) {
	return lw.loader.LoadDashboard(ctx, name)
}

// LoadBenchmark loads a benchmark and all its children on-demand.
func (lw *LazyWorkspace) LoadBenchmark(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
	return lw.loader.LoadBenchmark(ctx, name)
}

// LoadBenchmarkForExecution loads a benchmark with all children properly resolved
// and associated with their parents. This is needed for execution because the standard
// LoadBenchmark only caches resources but doesn't set the Children field on benchmarks.
//
// The key difference from LoadBenchmark:
// - LoadBenchmark: loads resources into cache, but GetChildren() returns empty
// - LoadBenchmarkForExecution: loads resources AND sets Children field properly
//
// Thread-safe: Uses per-benchmark sync.Once to ensure resolution happens only once
// even when called concurrently from multiple goroutines.
func (lw *LazyWorkspace) LoadBenchmarkForExecution(ctx context.Context, name string) (modconfig.ModTreeItem, error) {
	// First, load all resources into the cache using the standard loader
	benchmark, err := lw.loader.LoadBenchmark(ctx, name)
	if err != nil {
		return nil, err
	}

	// Get or create the resolution tracker for this benchmark
	// This ensures that child resolution only happens once per benchmark
	resolutionI, _ := lw.resolvedBenchmarks.LoadOrStore(name, &benchmarkResolution{})
	resolution := resolutionI.(*benchmarkResolution)

	// Resolve children only once (thread-safe)
	resolution.once.Do(func() {
		resolution.err = lw.resolveChildrenRecursively(ctx, benchmark)
	})

	if resolution.err != nil {
		return nil, resolution.err
	}

	return benchmark, nil
}

// resolveChildrenRecursively walks the benchmark tree and sets the Children field
// by resolving child names from the index and looking up resources from the cache.
func (lw *LazyWorkspace) resolveChildrenRecursively(ctx context.Context, item modconfig.ModTreeItem) error {
	// Get child names from the index
	entry, ok := lw.index.Get(item.Name())
	if !ok {
		return nil // No index entry, nothing to resolve
	}

	if len(entry.ChildNames) == 0 {
		return nil // No children to resolve
	}

	// Resolve each child from the cache and build the children slice
	children := make([]modconfig.ModTreeItem, 0, len(entry.ChildNames))
	for _, childName := range entry.ChildNames {
		// Load child from cache (should already be loaded by LoadBenchmark)
		childResource, err := lw.loader.Load(ctx, childName)
		if err != nil {
			return fmt.Errorf("failed to load child %s: %w", childName, err)
		}

		childItem, ok := childResource.(modconfig.ModTreeItem)
		if !ok {
			continue // Skip non-tree items
		}

		children = append(children, childItem)

		// Set parent relationship
		if err := childItem.AddParent(item); err != nil {
			return err
		}

		// Recursively resolve children for benchmarks (both regular and detection)
		switch childResource.(type) {
		case *resources.Benchmark, *resources.DetectionBenchmark:
			if err := lw.resolveChildrenRecursively(ctx, childItem); err != nil {
				return err
			}
		}
	}

	// Set children on the item using the ModTreeItemImpl's SetChildren method
	if impl := item.GetModTreeItemImpl(); impl != nil {
		impl.SetChildren(children)
		impl.ChildNameStrings = entry.ChildNames
	}

	return nil
}

// LoadResource loads a single resource by name.
func (lw *LazyWorkspace) LoadResource(ctx context.Context, name string) (modconfig.HclResource, error) {
	return lw.loader.Load(ctx, name)
}

// InvalidateResource removes a resource from the cache.
func (lw *LazyWorkspace) InvalidateResource(name string) {
	lw.cache.Invalidate(name)
}

// InvalidateAll clears the entire cache.
func (lw *LazyWorkspace) InvalidateAll() {
	lw.cache.Clear()
}

// CacheStats returns cache statistics.
func (lw *LazyWorkspace) CacheStats() resourcecache.CacheStats {
	return lw.cache.Stats()
}

// IndexStats returns index statistics.
func (lw *LazyWorkspace) IndexStats() resourceindex.IndexStats {
	return lw.index.Stats()
}

// Close cleans up the lazy workspace.
func (lw *LazyWorkspace) Close() {
	// Stop background resolution first
	lw.StopBackgroundResolution()

	lw.PowerpipeWorkspace.Close()
	lw.cache.Clear()
}

// PublishDashboardEvent publishes a dashboard event.
func (lw *LazyWorkspace) PublishDashboardEvent(ctx context.Context, event dashboardevents.DashboardEvent) {
	lw.PowerpipeWorkspace.PublishDashboardEvent(ctx, event)
}

// IsLazy returns true to indicate this is a lazy-loading workspace.
func (lw *LazyWorkspace) IsLazy() bool {
	return true
}

// StartBackgroundResolution begins background metadata resolution.
// This resolves variable references, templates, and function calls in the background,
// progressively updating the index as resolution completes.
func (lw *LazyWorkspace) StartBackgroundResolution() {
	if lw.backgroundResolver != nil {
		return // Already running
	}

	lw.backgroundResolver = NewBackgroundResolver(lw,
		WithWorkers(4),
		WithOnUpdate(lw.handleResourceUpdate),
		WithOnComplete(lw.handleResolutionComplete),
	)

	lw.backgroundResolver.Start()
}

// StopBackgroundResolution stops background resolution if running.
func (lw *LazyWorkspace) StopBackgroundResolution() {
	if lw.backgroundResolver != nil {
		lw.backgroundResolver.Stop()
		lw.backgroundResolver = nil
	}
}

// IsFullyResolved returns true if background resolution is complete.
func (lw *LazyWorkspace) IsFullyResolved() bool {
	return lw.fullyResolved
}

// handleResourceUpdate is called when a resource's metadata is resolved.
func (lw *LazyWorkspace) handleResourceUpdate(resourceName string) {
	lw.updateListenersMu.RLock()
	listeners := lw.updateListeners
	lw.updateListenersMu.RUnlock()

	for _, listener := range listeners {
		listener.OnResourceUpdated(resourceName)
	}
}

// handleResolutionComplete is called when all background resolution is done.
func (lw *LazyWorkspace) handleResolutionComplete() {
	slog.Info("Background resolution complete")
	lw.fullyResolved = true

	lw.updateListenersMu.RLock()
	listeners := lw.updateListeners
	lw.updateListenersMu.RUnlock()

	for _, listener := range listeners {
		listener.OnResolutionComplete()
	}
}

// RegisterUpdateListener adds a listener for background resolution updates.
func (lw *LazyWorkspace) RegisterUpdateListener(listener UpdateListener) {
	lw.updateListenersMu.Lock()
	defer lw.updateListenersMu.Unlock()
	lw.updateListeners = append(lw.updateListeners, listener)
}

// UnregisterUpdateListener removes an update listener.
func (lw *LazyWorkspace) UnregisterUpdateListener(listener UpdateListener) {
	lw.updateListenersMu.Lock()
	defer lw.updateListenersMu.Unlock()

	for i, l := range lw.updateListeners {
		if l == listener {
			lw.updateListeners = append(lw.updateListeners[:i], lw.updateListeners[i+1:]...)
			return
		}
	}
}

// PrioritizeResolution moves a resource to front of resolution queue.
// Useful when user is about to view a resource.
func (lw *LazyWorkspace) PrioritizeResolution(resourceName string) {
	if lw.backgroundResolver != nil {
		lw.backgroundResolver.Prioritize(resourceName)
	}
}

// ResolveNow immediately resolves a resource, bypassing the queue.
// Useful for on-demand resolution when user clicks a dashboard.
func (lw *LazyWorkspace) ResolveNow(ctx context.Context, resourceName string) error {
	if lw.backgroundResolver != nil {
		return lw.backgroundResolver.ResolveNow(ctx, resourceName)
	}

	// If no background resolver, create a temporary one just for this resolution
	entry, ok := lw.index.Get(resourceName)
	if !ok || entry.IsFullyResolved() {
		return nil
	}

	resolver := NewBackgroundResolver(lw)
	return resolver.ResolveNow(ctx, resourceName)
}

// BackgroundResolverStats returns statistics about background resolution.
func (lw *LazyWorkspace) BackgroundResolverStats() BackgroundResolverStats {
	if lw.backgroundResolver == nil {
		return BackgroundResolverStats{}
	}
	return lw.backgroundResolver.Stats()
}

// WaitForResolution waits for background resolution to complete, up to the specified timeout.
// Returns true if resolution completed, false if timeout was reached.
func (lw *LazyWorkspace) WaitForResolution(timeout time.Duration) bool {
	if lw.backgroundResolver == nil {
		return true
	}
	return lw.backgroundResolver.WaitForComplete(timeout)
}
