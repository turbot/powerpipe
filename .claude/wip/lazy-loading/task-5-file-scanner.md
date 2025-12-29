# Task 5: File Scanner for Index Building

## Objective

Implement a fast file scanner that extracts resource metadata from HCL files WITHOUT full parsing. This enables building the resource index with minimal memory and CPU usage.

## Context

- Full HCL parsing is expensive (creates AST, allocates memory)
- For the index, we only need: type, name, title, tags, parent/children, file location
- A lightweight scanner can extract this ~10x faster than full parsing
- This is where the startup speed improvement comes from

## Repository

**This task is Powerpipe-only.** No changes to pipe-fittings required.

The scanner is intentionally separate from pipe-fittings parsing because:
1. It's regex-based, not full HCL parsing
2. It's specific to Powerpipe resource types
3. It avoids coupling lazy loading to the shared library

## Dependencies

### Prerequisites
- Task 4 (Resource Index) - Need index structure to populate

### Files to Create (powerpipe)
- `internal/resourceindex/scanner.go` - HCL file scanner
- `internal/resourceindex/scanner_test.go` - Scanner tests

### Files to Modify
- None - scanner is standalone

## Implementation Details

### 1. Scanner Design

The scanner reads HCL files and extracts block headers without building full AST:

```go
// internal/resourceindex/scanner.go
package resourceindex

import (
    "bufio"
    "os"
    "regexp"
    "strings"
)

// ResourceTypes that should be indexed
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
    "locals":              true,
}

// Scanner extracts resource metadata from HCL files
type Scanner struct {
    modName string
    index   *ResourceIndex
}

func NewScanner(modName string) *Scanner {
    return &Scanner{
        modName: modName,
        index:   NewResourceIndex(),
    }
}

// ScanFile extracts index entries from a single HCL file
func (s *Scanner) ScanFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0
    var currentBlock *blockState

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // Check for block start
        if blockStart := s.parseBlockStart(line); blockStart != nil {
            // Close previous block if any
            if currentBlock != nil {
                s.finalizeBlock(currentBlock, lineNum-1, filePath)
            }

            currentBlock = &blockState{
                blockType:  blockStart.blockType,
                name:       blockStart.name,
                startLine:  lineNum,
                filePath:   filePath,
                attributes: make(map[string]string),
                children:   []string{},
            }
            continue
        }

        // Extract attributes from current block
        if currentBlock != nil {
            if attr := s.parseAttribute(line); attr != nil {
                currentBlock.attributes[attr.name] = attr.value
            }

            // Check for nested block (child)
            if nestedBlock := s.parseNestedBlockStart(line); nestedBlock != nil {
                currentBlock.children = append(currentBlock.children, nestedBlock)
            }

            // Check for block end
            if s.isBlockEnd(line) {
                s.finalizeBlock(currentBlock, lineNum, filePath)
                currentBlock = nil
            }
        }
    }

    // Handle unclosed block at EOF
    if currentBlock != nil {
        s.finalizeBlock(currentBlock, lineNum, filePath)
    }

    return scanner.Err()
}

type blockState struct {
    blockType  string
    name       string
    startLine  int
    filePath   string
    attributes map[string]string
    children   []string
}

type blockStart struct {
    blockType string
    name      string
}

type attribute struct {
    name  string
    value string
}

// Block start pattern: `blocktype "name" {` or `blocktype "name" "label" {`
var blockStartRegex = regexp.MustCompile(`^\s*(\w+)\s+"([^"]+)"(?:\s+"[^"]+")?\s*\{?\s*$`)

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

// Attribute pattern: `key = "value"` or `key = value`
var attrRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*"?([^"]*)"?\s*$`)

func (s *Scanner) parseAttribute(line string) *attribute {
    matches := attrRegex.FindStringSubmatch(line)
    if len(matches) < 3 {
        return nil
    }

    return &attribute{
        name:  matches[1],
        value: strings.TrimSpace(matches[2]),
    }
}

// Nested block pattern for children
var nestedBlockRegex = regexp.MustCompile(`^\s*(\w+)\s*=\s*(\w+)\.(\w+)\.(\w+)`)

func (s *Scanner) parseNestedBlockStart(line string) string {
    // Look for child references like: children = [benchmark.child1, benchmark.child2]
    // or: query = query.my_query
    matches := nestedBlockRegex.FindStringSubmatch(line)
    if len(matches) >= 5 {
        // Return full reference name
        return matches[2] + "." + matches[3] + "." + matches[4]
    }
    return ""
}

func (s *Scanner) isBlockEnd(line string) bool {
    trimmed := strings.TrimSpace(line)
    return trimmed == "}"
}

func (s *Scanner) finalizeBlock(block *blockState, endLine int, filePath string) {
    fullName := s.modName + "." + block.blockType + "." + block.name

    entry := &IndexEntry{
        Type:      block.blockType,
        Name:      fullName,
        ShortName: block.name,
        FileName:  filePath,
        StartLine: block.startLine,
        EndLine:   endLine,
        ModName:   s.modName,
    }

    // Extract common attributes
    if title, ok := block.attributes["title"]; ok {
        entry.Title = title
    }
    if desc, ok := block.attributes["description"]; ok {
        entry.Description = desc
    }
    if sql, ok := block.attributes["sql"]; ok {
        entry.HasSQL = sql != ""
    }

    // Extract tags (simplified - full parsing needed for complex tags)
    // Tags are often in format: tags = { key = "value" }
    // For now, we'll do basic extraction

    entry.ChildNames = block.children

    s.index.Add(entry)
}

// ScanDirectory scans all .pp files in a directory
func (s *Scanner) ScanDirectory(dirPath string) error {
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() {
            // Recursively scan subdirectories
            if err := s.ScanDirectory(filepath.Join(dirPath, entry.Name())); err != nil {
                return err
            }
            continue
        }

        // Only scan .pp files
        if !strings.HasSuffix(entry.Name(), ".pp") {
            continue
        }

        filePath := filepath.Join(dirPath, entry.Name())
        if err := s.ScanFile(filePath); err != nil {
            return fmt.Errorf("scanning %s: %w", filePath, err)
        }
    }

    return nil
}

// GetIndex returns the built index
func (s *Scanner) GetIndex() *ResourceIndex {
    return s.index
}
```

### 2. Enhanced Scanner with Byte Offsets

```go
// ScanFileWithOffsets extracts entries with byte offsets for efficient seeking
func (s *Scanner) ScanFileWithOffsets(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    lineNum := 0
    byteOffset := int64(0)
    var currentBlock *blockStateWithOffset

    for {
        line, err := reader.ReadString('\n')
        if err != nil && err != io.EOF {
            return err
        }

        lineLen := int64(len(line))
        lineNum++

        // Check for block start
        if blockStart := s.parseBlockStart(line); blockStart != nil {
            if currentBlock != nil {
                s.finalizeBlockWithOffset(currentBlock, lineNum-1,
                    byteOffset-lineLen, filePath)
            }

            currentBlock = &blockStateWithOffset{
                blockState: blockState{
                    blockType:  blockStart.blockType,
                    name:       blockStart.name,
                    startLine:  lineNum,
                    filePath:   filePath,
                    attributes: make(map[string]string),
                },
                startOffset: byteOffset - lineLen,
            }
        }

        // ... rest of parsing logic

        byteOffset += lineLen

        if err == io.EOF {
            break
        }
    }

    return nil
}

type blockStateWithOffset struct {
    blockState
    startOffset int64
}
```

### 3. Parallel Directory Scanning

```go
// ScanDirectoryParallel scans files in parallel for faster indexing
func (s *Scanner) ScanDirectoryParallel(dirPath string, workers int) error {
    // Collect all .pp files
    var files []string
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasSuffix(path, ".pp") {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    // Process in parallel
    if workers <= 0 {
        workers = runtime.NumCPU()
    }

    fileChan := make(chan string, len(files))
    errChan := make(chan error, workers)
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            localIndex := NewResourceIndex()

            for filePath := range fileChan {
                localScanner := &Scanner{
                    modName: s.modName,
                    index:   localIndex,
                }
                if err := localScanner.ScanFile(filePath); err != nil {
                    errChan <- err
                    return
                }
            }

            // Merge local index into main index
            s.mergeIndex(localIndex)
        }()
    }

    // Send files to workers
    for _, f := range files {
        fileChan <- f
    }
    close(fileChan)

    wg.Wait()
    close(errChan)

    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *Scanner) mergeIndex(other *ResourceIndex) {
    s.index.mu.Lock()
    defer s.index.mu.Unlock()

    other.mu.RLock()
    defer other.mu.RUnlock()

    for name, entry := range other.entries {
        s.index.entries[name] = entry
        if s.index.byType[entry.Type] == nil {
            s.index.byType[entry.Type] = make(map[string]*IndexEntry)
        }
        s.index.byType[entry.Type][name] = entry
        s.index.totalSize += entry.Size()
    }
}
```

### 4. Scanner Tests

```go
// internal/resourceindex/scanner_test.go
package resourceindex

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestScanner_SimpleFile(t *testing.T) {
    // Create temp file
    tmpDir := t.TempDir()
    filePath := filepath.Join(tmpDir, "test.pp")

    content := `
dashboard "my_dashboard" {
    title = "My Dashboard"
    description = "A test dashboard"

    container {
        card {
            sql = "SELECT 1"
        }
    }
}

query "my_query" {
    title = "My Query"
    sql = "SELECT * FROM table"
}
`
    require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

    scanner := NewScanner("testmod")
    err := scanner.ScanFile(filePath)
    require.NoError(t, err)

    index := scanner.GetIndex()

    // Check dashboard was found
    dash, ok := index.Get("testmod.dashboard.my_dashboard")
    assert.True(t, ok)
    assert.Equal(t, "My Dashboard", dash.Title)
    assert.Equal(t, "A test dashboard", dash.Description)

    // Check query was found
    query, ok := index.Get("testmod.query.my_query")
    assert.True(t, ok)
    assert.Equal(t, "My Query", query.Title)
    assert.True(t, query.HasSQL)
}

func TestScanner_BenchmarkHierarchy(t *testing.T) {
    tmpDir := t.TempDir()
    filePath := filepath.Join(tmpDir, "benchmarks.pp")

    content := `
benchmark "parent" {
    title = "Parent Benchmark"
    children = [
        benchmark.child1,
        benchmark.child2
    ]
}

benchmark "child1" {
    title = "Child 1"
}

benchmark "child2" {
    title = "Child 2"
}
`
    require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

    scanner := NewScanner("testmod")
    err := scanner.ScanFile(filePath)
    require.NoError(t, err)

    index := scanner.GetIndex()

    parent, ok := index.Get("testmod.benchmark.parent")
    assert.True(t, ok)
    assert.Contains(t, parent.ChildNames, "benchmark.child1")
    assert.Contains(t, parent.ChildNames, "benchmark.child2")
}

func TestScanner_Performance(t *testing.T) {
    // Create a large mod directory
    tmpDir := t.TempDir()

    // Create 100 files with 10 resources each
    for i := 0; i < 100; i++ {
        filePath := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pp", i))
        var content strings.Builder
        for j := 0; j < 10; j++ {
            content.WriteString(fmt.Sprintf(`
query "query_%d_%d" {
    title = "Query %d %d"
    sql = "SELECT %d"
}
`, i, j, i, j, j))
        }
        require.NoError(t, os.WriteFile(filePath, []byte(content.String()), 0644))
    }

    // Time the scan
    start := time.Now()

    scanner := NewScanner("testmod")
    err := scanner.ScanDirectory(tmpDir)
    require.NoError(t, err)

    elapsed := time.Since(start)
    index := scanner.GetIndex()

    t.Logf("Scanned %d resources in %v", index.Count(), elapsed)
    t.Logf("Index size: %d bytes", index.Size())

    assert.Equal(t, 1000, index.Count())
    assert.Less(t, elapsed.Milliseconds(), int64(500), "Scan too slow")
}
```

## Acceptance Criteria

- [ ] Scanner extracts resource type, name, title, description
- [ ] Scanner captures file path, start/end lines, byte offsets
- [ ] Scanner identifies parent-child relationships
- [ ] Scanner handles all 20+ resource types
- [ ] Scanner processes 1000 resources in < 100ms
- [ ] Parallel scanning available for large mods
- [ ] Scanner handles malformed files gracefully
- [ ] All unit tests pass
- [ ] Scanner output matches full HCL parse for metadata

## Notes

- Scanner is intentionally simple - regex-based, not full parser
- May need to handle edge cases (multi-line strings, heredocs)
- Consider falling back to full parse for complex files
- Byte offsets enable efficient seeking in Task 7
