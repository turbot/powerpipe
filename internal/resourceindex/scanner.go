package resourceindex

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
// This is much faster than full HCL parsing and creates minimal allocations.
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

// blockState tracks state while parsing a block
type blockState struct {
	blockType   string
	name        string
	startLine   int
	startOffset int64
	filePath    string
	depth       int // brace nesting depth
	attributes  map[string]string
	children    []string
	tags        map[string]string
	inTags      bool // tracking if we're inside a tags block
	inChildren  bool // tracking if we're inside a children array
}

// blockStart contains parsed block start info
type blockStart struct {
	blockType string
	name      string
}

// attribute contains parsed attribute info
type attribute struct {
	name  string
	value string
}

// Block start pattern: `blocktype "name" {` or `blocktype "name" "label" {`
// Also handles blocks without opening brace on same line, or single-line blocks
var blockStartRegex = regexp.MustCompile(`^\s*(\w+)\s+"([^"]+)"(?:\s+"[^"]*")?\s*\{?`)

// Attribute pattern: `key = "value"` (quoted string)
var attrStringRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*"([^"]*)"`)

// Attribute pattern: `key = value` (unquoted, for booleans, numbers, references)
var attrUnquotedRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*([^\s"#]+)`)

// Child reference pattern: `type.name` (e.g., control.my_control, benchmark.child)
var childRefRegex = regexp.MustCompile(`\b(\w+)\.(\w+)\b`)

// Tag entry pattern: `key = "value"` inside tags block
var tagEntryRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*"([^"]*)"`)

// ScanFile extracts index entries from a single HCL file.
// It uses a line-by-line scanner for efficiency.
func (s *Scanner) ScanFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.scanReader(file, filePath)
}

// ScanFileWithOffsets extracts entries with byte offsets for efficient seeking.
func (s *Scanner) ScanFileWithOffsets(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.scanReaderWithOffsets(file, filePath)
}

func (s *Scanner) scanReader(r io.Reader, filePath string) error {
	scanner := bufio.NewScanner(r)
	// Increase buffer size for files with long lines (heredocs)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	var blockStack []*blockState

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for block start
		if blockStart := s.parseBlockStart(line); blockStart != nil {
			block := &blockState{
				blockType:  blockStart.blockType,
				name:       blockStart.name,
				startLine:  lineNum,
				filePath:   filePath,
				depth:      0,
				attributes: make(map[string]string),
				children:   []string{},
				tags:       make(map[string]string),
			}

			// Count braces on this line
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			block.depth = openBraces - closeBraces

			// Process any attributes on the same line as block start
			// (for single-line blocks like: dashboard "d1" { title = "Test" })
			if openBraces > 0 {
				s.processBlockLine(line, block)
			}

			// If block is already closed on this line
			if block.depth <= 0 && openBraces > 0 {
				s.finalizeBlock(block, lineNum, filePath, 0, 0)
			} else {
				blockStack = append(blockStack, block)
			}
			continue
		}

		// Process current block(s)
		if len(blockStack) > 0 {
			currentBlock := blockStack[len(blockStack)-1]
			s.processBlockLine(line, currentBlock)

			// Check for block end
			if currentBlock.depth <= 0 {
				// Finalize and pop the block
				s.finalizeBlock(currentBlock, lineNum, filePath, 0, 0)
				blockStack = blockStack[:len(blockStack)-1]
			}
		}
	}

	// Handle any unclosed blocks
	for i := len(blockStack) - 1; i >= 0; i-- {
		s.finalizeBlock(blockStack[i], lineNum, filePath, 0, 0)
	}

	return scanner.Err()
}

func (s *Scanner) scanReaderWithOffsets(r io.Reader, filePath string) error {
	reader := bufio.NewReader(r)
	lineNum := 0
	byteOffset := int64(0)
	var blockStack []*blockState

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lineLen := int64(len(line))
		lineNum++
		lineStart := byteOffset

		// Check for block start
		if blockStart := s.parseBlockStart(line); blockStart != nil {
			block := &blockState{
				blockType:   blockStart.blockType,
				name:        blockStart.name,
				startLine:   lineNum,
				startOffset: lineStart,
				filePath:    filePath,
				depth:       0,
				attributes:  make(map[string]string),
				children:    []string{},
				tags:        make(map[string]string),
			}

			// Count braces on this line
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			block.depth = openBraces - closeBraces

			// Process any attributes on the same line as block start
			if openBraces > 0 {
				s.processBlockLine(line, block)
			}

			// If block is already closed on this line
			if block.depth <= 0 && openBraces > 0 {
				endOffset := byteOffset + lineLen
				byteLen := int(endOffset - lineStart)
				s.finalizeBlock(block, lineNum, filePath, lineStart, byteLen)
			} else {
				blockStack = append(blockStack, block)
			}
		} else if len(blockStack) > 0 {
			currentBlock := blockStack[len(blockStack)-1]
			s.processBlockLine(line, currentBlock)

			// Check for block end
			if currentBlock.depth <= 0 {
				endOffset := byteOffset + lineLen
				byteLen := int(endOffset - currentBlock.startOffset)
				s.finalizeBlock(currentBlock, lineNum, filePath, currentBlock.startOffset, byteLen)
				blockStack = blockStack[:len(blockStack)-1]
			}
		}

		byteOffset += lineLen

		if err == io.EOF {
			break
		}
	}

	// Handle any unclosed blocks
	for i := len(blockStack) - 1; i >= 0; i-- {
		block := blockStack[i]
		byteLen := int(byteOffset - block.startOffset)
		s.finalizeBlock(block, lineNum, filePath, block.startOffset, byteLen)
	}

	return nil
}

func (s *Scanner) processBlockLine(line string, block *blockState) {
	trimmed := strings.TrimSpace(line)

	// Track brace depth (but don't double-count on block start line)
	openBraces := strings.Count(line, "{")
	closeBraces := strings.Count(line, "}")

	// Only skip brace counting for lines that started a NEW indexed block.
	// Lines that match the regex pattern but aren't indexed types (like "column")
	// should still have their braces counted.
	skipBraceCount := false
	if matches := blockStartRegex.FindStringSubmatch(line); len(matches) >= 3 {
		blockType := matches[1]
		// Only skip if this is an indexed type (which means we created a block for it)
		if indexedTypes[blockType] {
			skipBraceCount = true
		}
	}

	if !skipBraceCount {
		block.depth += openBraces - closeBraces
	}

	// Check for tags block start
	if strings.HasPrefix(trimmed, "tags") && strings.Contains(line, "{") {
		block.inTags = true
		return
	}

	// Check for children array start
	if strings.HasPrefix(trimmed, "children") && strings.Contains(line, "[") {
		block.inChildren = true
	}

	// Parse tags entries
	if block.inTags {
		if matches := tagEntryRegex.FindStringSubmatch(line); len(matches) >= 3 {
			block.tags[matches[1]] = matches[2]
		}
		if strings.Contains(line, "}") && !strings.Contains(line, "{") {
			block.inTags = false
		}
		return
	}

	// Parse children references
	if block.inChildren {
		// Look for type.name references
		matches := childRefRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				childType := match[1]
				childName := match[2]
				// Build full child name: mod.type.name
				fullChildName := s.modName + "." + childType + "." + childName
				block.children = append(block.children, fullChildName)
			}
		}
		if strings.Contains(line, "]") {
			block.inChildren = false
		}
		return
	}

	// Parse standard attributes
	if attr := s.parseAttribute(line); attr != nil {
		block.attributes[attr.name] = attr.value
	}

	// Check for query reference: query = type.name
	if strings.Contains(trimmed, "query") && strings.Contains(trimmed, "=") {
		if matches := childRefRegex.FindStringSubmatch(line); len(matches) >= 3 {
			// This is a query reference, not inline SQL
			block.attributes["query_ref"] = matches[1] + "." + matches[2]
		}
	}
}

func (s *Scanner) parseBlockStart(line string) *blockStart {
	matches := blockStartRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}

	blockType := matches[1]
	name := matches[2]

	if !indexedTypes[blockType] {
		return nil
	}

	return &blockStart{
		blockType: blockType,
		name:      name,
	}
}

func (s *Scanner) parseAttribute(line string) *attribute {
	// Try quoted string first
	if matches := attrStringRegex.FindStringSubmatch(line); len(matches) >= 3 {
		return &attribute{
			name:  matches[1],
			value: matches[2],
		}
	}

	// Try unquoted value
	if matches := attrUnquotedRegex.FindStringSubmatch(line); len(matches) >= 3 {
		return &attribute{
			name:  matches[1],
			value: strings.TrimSpace(matches[2]),
		}
	}

	return nil
}

func (s *Scanner) finalizeBlock(block *blockState, endLine int, filePath string, byteOffset int64, byteLen int) {
	// Use interned strings for commonly repeated values
	internedModName := intern.Intern(s.modName)
	internedType := intern.Intern(block.blockType)
	internedShortName := intern.Intern(block.name)
	internedFileName := intern.Intern(filePath)

	// Build full name with interned components
	fullName := intern.Intern(internedModName + "." + internedType + "." + internedShortName)

	// Intern mod root if set (paths repeat across files)
	internedModRoot := ""
	if s.modRoot != "" {
		internedModRoot = intern.Intern(s.modRoot)
	}

	entry := &IndexEntry{
		Type:       internedType,
		Name:       fullName,
		ShortName:  internedShortName,
		FileName:   internedFileName,
		StartLine:  block.startLine,
		EndLine:    endLine,
		ByteOffset: byteOffset,
		ByteLength: byteLen,
		ModName:    internedModName,
		ModRoot:    internedModRoot,
	}

	// Extract common attributes
	// Titles and descriptions are usually unique, don't intern them
	if title, ok := block.attributes["title"]; ok {
		entry.Title = title
	}
	if desc, ok := block.attributes["description"]; ok {
		entry.Description = desc
	}

	// Check for SQL (inline or referenced)
	if _, ok := block.attributes["sql"]; ok {
		entry.HasSQL = true
	}
	if queryRef, ok := block.attributes["query_ref"]; ok {
		entry.HasSQL = true // References a query with SQL
		// Build full query reference: mod.type.name (interned)
		entry.QueryRef = intern.Intern(internedModName + "." + queryRef)
	}

	// Copy tags with interned keys (tag keys are often repeated)
	if len(block.tags) > 0 {
		entry.Tags = make(map[string]string, len(block.tags))
		for k, v := range block.tags {
			// Intern keys (often repeated: service, category, etc.)
			// Values are usually unique, don't intern
			entry.Tags[intern.Intern(k)] = v
		}
	}

	// Set benchmark type (using interned constant strings)
	if block.blockType == "benchmark" {
		entry.BenchmarkType = intern.BenchmarkTypeControl
	} else if block.blockType == "detection_benchmark" {
		entry.BenchmarkType = intern.BenchmarkTypeDetection
	}

	// Copy children with interned names
	if len(block.children) > 0 {
		entry.ChildNames = intern.InternSlice(block.children)
	}

	s.addEntry(entry)
}

func (s *Scanner) addEntry(entry *IndexEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index.Add(entry)
}

// ScanDirectory scans all .pp files in a directory recursively.
func (s *Scanner) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Only scan .pp files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".pp") {
			return nil
		}

		return s.ScanFile(path)
	})
}

// ScanDirectoryParallel scans files in parallel for faster indexing.
func (s *Scanner) ScanDirectoryParallel(dirPath string, workers int) error {
	// Collect all .pp files
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(path, ".pp") {
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

		// Only scan .pp files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".pp") {
			return nil
		}

		return s.ScanFile(path)
	})

	// Restore original mod name and root
	s.modName = originalModName
	s.modRoot = originalModRoot

	return err
}

// ScanBytes scans HCL content from a byte slice.
// Useful for testing without file I/O.
func (s *Scanner) ScanBytes(content []byte, filePath string) error {
	return s.scanReader(strings.NewReader(string(content)), filePath)
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
