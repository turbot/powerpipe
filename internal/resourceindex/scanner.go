package resourceindex

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/turbot/powerpipe/internal/intern"
)

// ResourceTypes that should be indexed for lazy loading
var indexedTypes = map[string]bool{
	"dashboard":           true,
	"benchmark":           true,
	"control":             true,
	"query":               true,
	"card":                true,
	"chart":               true,
	"container":           true,
	"flow":                true,
	"graph":               true,
	"hierarchy":           true,
	"image":               true,
	"input":               true,
	"node":                true,
	"edge":                true,
	"table":               true,
	"text":                true,
	"category":            true,
	"detection":           true,
	"detection_benchmark": true,
	"variable":            true,
	"with":                true,
}

// Scanner extracts resource metadata from HCL files without full parsing.
// This uses the HCL syntax parser for correctness while avoiding full expression evaluation.
type Scanner struct {
	modName string
	modRoot string // Root directory of the current mod (for file() function resolution)
	index   *ResourceIndex
	mu      sync.Mutex // protects index during parallel scanning
}

// NewScanner creates a new scanner for extracting resource metadata.
func NewScanner(modName string) *Scanner {
	return &Scanner{
		modName: modName,
		index:   NewResourceIndex(),
	}
}

// SetModRoot sets the root directory for the current mod.
func (s *Scanner) SetModRoot(root string) {
	s.modRoot = root
}

// ScanFile extracts index entries from a single HCL file.
// It uses the HCL syntax parser for correctness.
func (s *Scanner) ScanFile(filePath string) error {
	return s.ScanFileHCL(filePath)
}

// ScanFileWithOffsets extracts entries with byte offsets for efficient seeking.
func (s *Scanner) ScanFileWithOffsets(filePath string) error {
	return s.ScanFileHCLWithOffsets(filePath)
}

// ScanBytes scans HCL content from a byte slice.
// Useful for testing without file I/O.
func (s *Scanner) ScanBytes(content []byte, filePath string) error {
	return s.ScanBytesHCL(content, filePath)
}

func (s *Scanner) addEntry(entry *IndexEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index.Add(entry)
}

// isPowerpipeFile checks if a file has a valid Powerpipe extension (.pp or .sp)
func isPowerpipeFile(name string) bool {
	return strings.HasSuffix(name, ".pp") || strings.HasSuffix(name, ".sp")
}

// ScanDirectory scans all .pp and .sp files in a directory recursively.
func (s *Scanner) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Only scan .pp and .sp files
		if info.IsDir() || !isPowerpipeFile(info.Name()) {
			return nil
		}

		return s.ScanFile(path)
	})
}

// ScanDirectoryParallel scans files in parallel for faster indexing.
func (s *Scanner) ScanDirectoryParallel(dirPath string, workers int) error {
	// Collect all .pp and .sp files
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !info.IsDir() && isPowerpipeFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	// Set default workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if workers > len(files) {
		workers = len(files)
	}

	// Process in parallel
	fileChan := make(chan string, len(files))
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				if err := s.ScanFile(filePath); err != nil {
					select {
					case errChan <- fmt.Errorf("scanning %s: %w", filePath, err):
					default:
					}
					return
				}
			}
		}()
	}

	// Send files to workers
	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	// Wait for completion
	wg.Wait()
	close(errChan)

	// Return first error if any
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// ScanDirectoryWithModName scans a directory with a specific mod name.
// This is used for scanning dependency mods where each mod has its own name.
func (s *Scanner) ScanDirectoryWithModName(dirPath, modName string) error {
	// Save current mod name and root
	originalModName := s.modName
	originalModRoot := s.modRoot

	// Set the new mod name and root for this directory
	s.modName = modName
	s.modRoot = dirPath

	// Scan the directory (don't recurse into subdirectories with mod.pp)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Only scan .pp and .sp files
		if info.IsDir() || !isPowerpipeFile(info.Name()) {
			return nil
		}

		return s.ScanFile(path)
	})

	// Restore original mod name and root
	s.modName = originalModName
	s.modRoot = originalModRoot

	return err
}

// GetIndex returns the built index.
func (s *Scanner) GetIndex() *ResourceIndex {
	return s.index
}

// SetModInfo sets mod-level information on the index.
func (s *Scanner) SetModInfo(modName, modFullName, modTitle string) {
	// Intern mod names as they're repeated across all resources
	s.index.ModName = intern.Intern(modName)
	s.index.ModFullName = intern.Intern(modFullName)
	// Titles are usually unique, don't intern
	s.index.ModTitle = modTitle
}

// MarkTopLevelResources marks resources as top-level based on parent references.
// This should be called after scanning is complete.
func (s *Scanner) MarkTopLevelResources() {
	s.index.mu.Lock()
	defer s.index.mu.Unlock()

	// Build set of all child names
	childNames := make(map[string]bool)
	for _, entry := range s.index.entries {
		for _, child := range entry.ChildNames {
			childNames[child] = true
		}
	}

	// Mark entries that are not children of anything as top-level
	for _, entry := range s.index.entries {
		// Only mark benchmarks and dashboards as top-level
		if entry.Type == "benchmark" || entry.Type == "detection_benchmark" || entry.Type == "dashboard" {
			if !childNames[entry.Name] {
				entry.IsTopLevel = true
			}
		}
	}
}

// SetParentNames sets ParentName on child entries based on ChildNames.
// This should be called after scanning is complete.
func (s *Scanner) SetParentNames() {
	s.index.mu.Lock()
	defer s.index.mu.Unlock()

	for _, entry := range s.index.entries {
		for _, childName := range entry.ChildNames {
			if child, ok := s.index.entries[childName]; ok {
				// Parent name is already interned (it's entry.Name)
				child.ParentName = entry.Name
			}
		}
	}
}
