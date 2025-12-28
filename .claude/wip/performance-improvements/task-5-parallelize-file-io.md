# Task 5: Parallelize File I/O

## Objective

Parallelize file reading in the mod loading process to reduce I/O wait time when loading many `.pp` files.

## Context

- Current `LoadFileData()` in pipe-fittings reads files sequentially
- For mods with 100+ files, sequential I/O is a significant bottleneck
- Modern SSDs can handle parallel reads efficiently
- This requires changes to the pipe-fittings library

## Dependencies

### Prerequisites
- Task 4 (Baseline Measurement) - Need baseline data for comparison

### Files to Modify (in pipe-fittings)
- `parse/parser.go` - `LoadFileData()` function

### Related Files in Powerpipe
- None directly modified, but will benefit from pipe-fittings change

## Implementation Details

### 1. Current Implementation (pipe-fittings)

```go
// Current: Sequential file reads
func LoadFileData(paths ...string) (map[string][]byte, hcl.Diagnostics) {
    var diags hcl.Diagnostics
    var fileData = map[string][]byte{}

    for _, configPath := range paths {
        data, err := os.ReadFile(configPath)  // Sequential!
        if err != nil {
            diags = append(diags, &hcl.Diagnostic{...})
            continue
        }
        fileData[configPath] = data
    }
    return fileData, diags
}
```

### 2. Optimized Implementation

```go
package parse

import (
    "os"
    "runtime"
    "sync"

    "github.com/hashicorp/hcl/v2"
)

// LoadFileData reads multiple files in parallel
func LoadFileData(paths ...string) (map[string][]byte, hcl.Diagnostics) {
    if len(paths) == 0 {
        return map[string][]byte{}, nil
    }

    // For small number of files, sequential is fine
    if len(paths) < 4 {
        return loadFileDataSequential(paths)
    }

    return loadFileDataParallel(paths)
}

func loadFileDataSequential(paths []string) (map[string][]byte, hcl.Diagnostics) {
    var diags hcl.Diagnostics
    fileData := make(map[string][]byte, len(paths))

    for _, configPath := range paths {
        data, err := os.ReadFile(configPath)
        if err != nil {
            diags = append(diags, &hcl.Diagnostic{
                Severity: hcl.DiagWarning,
                Summary:  fmt.Sprintf("failed to read config file %s", configPath),
                Detail:   err.Error(),
            })
            continue
        }
        fileData[configPath] = data
    }
    return fileData, diags
}

type fileReadResult struct {
    path string
    data []byte
    err  error
}

func loadFileDataParallel(paths []string) (map[string][]byte, hcl.Diagnostics) {
    var diags hcl.Diagnostics
    fileData := make(map[string][]byte, len(paths))

    // Use worker pool pattern
    numWorkers := runtime.NumCPU()
    if numWorkers > 8 {
        numWorkers = 8 // Cap at 8 to avoid too many open files
    }
    if numWorkers > len(paths) {
        numWorkers = len(paths)
    }

    pathsChan := make(chan string, len(paths))
    resultsChan := make(chan fileReadResult, len(paths))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range pathsChan {
                data, err := os.ReadFile(path)
                resultsChan <- fileReadResult{path: path, data: data, err: err}
            }
        }()
    }

    // Send work
    for _, path := range paths {
        pathsChan <- path
    }
    close(pathsChan)

    // Wait for completion in separate goroutine
    go func() {
        wg.Wait()
        close(resultsChan)
    }()

    // Collect results
    for result := range resultsChan {
        if result.err != nil {
            diags = append(diags, &hcl.Diagnostic{
                Severity: hcl.DiagWarning,
                Summary:  fmt.Sprintf("failed to read config file %s", result.path),
                Detail:   result.err.Error(),
            })
            continue
        }
        fileData[result.path] = result.data
    }

    return fileData, diags
}
```

### 3. Add Unit Tests

```go
// parse/parser_test.go

func TestLoadFileDataParallel(t *testing.T) {
    // Create temp directory with multiple files
    tmpDir := t.TempDir()

    numFiles := 50
    expectedData := make(map[string]string)

    for i := 0; i < numFiles; i++ {
        path := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pp", i))
        content := fmt.Sprintf("query \"q%d\" { sql = \"SELECT %d\" }", i, i)
        os.WriteFile(path, []byte(content), 0644)
        expectedData[path] = content
    }

    paths := make([]string, 0, numFiles)
    for path := range expectedData {
        paths = append(paths, path)
    }

    // Test parallel loading
    fileData, diags := LoadFileData(paths...)

    assert.Empty(t, diags)
    assert.Len(t, fileData, numFiles)

    for path, expected := range expectedData {
        assert.Equal(t, expected, string(fileData[path]))
    }
}

func TestLoadFileDataSequentialForSmallSets(t *testing.T) {
    // Verify small file sets still work correctly
    tmpDir := t.TempDir()

    paths := []string{
        filepath.Join(tmpDir, "a.pp"),
        filepath.Join(tmpDir, "b.pp"),
    }

    for _, p := range paths {
        os.WriteFile(p, []byte("test"), 0644)
    }

    fileData, diags := LoadFileData(paths...)
    assert.Empty(t, diags)
    assert.Len(t, fileData, 2)
}

func TestLoadFileDataHandlesMissingFiles(t *testing.T) {
    paths := []string{"/nonexistent/file.pp"}

    fileData, diags := LoadFileData(paths...)

    assert.Len(t, diags, 1)
    assert.Len(t, fileData, 0)
}
```

### 4. Benchmark the Change

```go
func BenchmarkLoadFileData_Sequential(b *testing.B) {
    paths := setupBenchmarkFiles(b, 100)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        loadFileDataSequential(paths)
    }
}

func BenchmarkLoadFileData_Parallel(b *testing.B) {
    paths := setupBenchmarkFiles(b, 100)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        loadFileDataParallel(paths)
    }
}

func setupBenchmarkFiles(b *testing.B, count int) []string {
    b.Helper()
    tmpDir := b.TempDir()
    paths := make([]string, count)

    for i := 0; i < count; i++ {
        path := filepath.Join(tmpDir, fmt.Sprintf("file_%d.pp", i))
        // Create realistic-sized file content
        content := strings.Repeat(fmt.Sprintf("query \"q%d\" { sql = \"SELECT %d\" }\n", i, i), 10)
        os.WriteFile(path, []byte(content), 0644)
        paths[i] = path
    }

    return paths
}
```

### 5. Measure Performance Improvement

After implementing, run:

```bash
# In powerpipe directory
POWERPIPE_TIMING=detailed go test -bench=BenchmarkLoadWorkspace -benchmem \
    ./internal/workspace/... -run=^$ \
    | tee benchmark_results/after_parallel_io.txt

# Compare with baseline
go run scripts/compare_benchmarks.go \
    benchmark_results/baseline/workspace_load.json \
    benchmark_results/after_parallel_io.json
```

## Acceptance Criteria

- [ ] `LoadFileData()` uses parallel reads for 4+ files
- [ ] Sequential fallback for small file sets
- [ ] Worker pool limited to avoid too many open files
- [ ] Unit tests pass for parallel loading
- [ ] Unit tests pass for edge cases (missing files, empty paths)
- [ ] No race conditions (verify with `go test -race`)
- [ ] Benchmark shows improvement for large file sets
- [ ] Performance results documented
- [ ] Changes compatible with existing API

## Expected Performance Improvement

| Mod Size | Files | Baseline | After | Improvement |
|----------|-------|----------|-------|-------------|
| Small | 5 | ~5ms | ~5ms | 0% (no change expected) |
| Medium | 5 | ~15ms | ~15ms | 0% |
| Large (many files) | 50+ | ~50ms | ~15ms | ~70% |

## Notes

- This change is in pipe-fittings, requires separate PR
- Consider file descriptor limits on different OSes
- May need to adjust worker count based on benchmarks
- Error handling must remain consistent
- Diagnostics order may change (document this)
