# Task 14: String Interning

## Objective

Implement string interning to deduplicate repeated strings across resources, reducing memory usage for common values like resource types, mod names, and tags.

## Context

- Many strings are repeated across resources (type names, mod name, common tags)
- Each duplicate is a separate allocation
- String interning shares single instances
- Can reduce memory for string-heavy workloads
- Lower priority optimization - implement if time permits

## Dependencies

### Prerequisites
- Task 9 (Workspace Integration) - Lazy workspace working
- Task 12 (Post-Parse Cleanup) - Cleanup infrastructure

### Files to Create
- `internal/intern/string_intern.go` - String interning implementation
- `internal/intern/string_intern_test.go` - Tests

### Files to Modify
- `internal/resourceindex/entry.go` - Use interned strings
- Various resource creation code paths

## Implementation Details

### 1. String Interner

```go
// internal/intern/string_intern.go
package intern

import (
    "sync"
)

// StringInterner deduplicates strings to reduce memory
type StringInterner struct {
    mu      sync.RWMutex
    strings map[string]string
    hits    int64
    misses  int64
}

// NewStringInterner creates a new interner
func NewStringInterner() *StringInterner {
    return &StringInterner{
        strings: make(map[string]string),
    }
}

// DefaultInterner is the global string interner
var DefaultInterner = NewStringInterner()

// Intern returns a canonical version of the string
// If the string was seen before, returns the original instance
func (i *StringInterner) Intern(s string) string {
    if s == "" {
        return ""
    }

    // Fast path: check if already interned
    i.mu.RLock()
    if interned, ok := i.strings[s]; ok {
        i.hits++
        i.mu.RUnlock()
        return interned
    }
    i.mu.RUnlock()

    // Slow path: add to intern table
    i.mu.Lock()
    defer i.mu.Unlock()

    // Double-check after acquiring write lock
    if interned, ok := i.strings[s]; ok {
        i.hits++
        return interned
    }

    // Intern the string
    i.strings[s] = s
    i.misses++
    return s
}

// InternSlice interns all strings in a slice
func (i *StringInterner) InternSlice(ss []string) []string {
    if len(ss) == 0 {
        return ss
    }

    result := make([]string, len(ss))
    for j, s := range ss {
        result[j] = i.Intern(s)
    }
    return result
}

// InternMap interns all keys and values in a map
func (i *StringInterner) InternMap(m map[string]string) map[string]string {
    if len(m) == 0 {
        return m
    }

    result := make(map[string]string, len(m))
    for k, v := range m {
        result[i.Intern(k)] = i.Intern(v)
    }
    return result
}

// Stats returns interning statistics
func (i *StringInterner) Stats() InternStats {
    i.mu.RLock()
    defer i.mu.RUnlock()

    // Calculate memory saved
    var savedBytes int64
    for s := range i.strings {
        // Each hit saves len(s) bytes (approximate)
        savedBytes += int64(len(s)) * i.hits / int64(len(i.strings))
    }

    return InternStats{
        UniqueStrings: len(i.strings),
        Hits:          i.hits,
        Misses:        i.misses,
        HitRate:       float64(i.hits) / float64(i.hits+i.misses),
        SavedBytes:    savedBytes,
    }
}

// Clear clears the intern table
func (i *StringInterner) Clear() {
    i.mu.Lock()
    defer i.mu.Unlock()

    i.strings = make(map[string]string)
    i.hits = 0
    i.misses = 0
}

type InternStats struct {
    UniqueStrings int
    Hits          int64
    Misses        int64
    HitRate       float64
    SavedBytes    int64
}

// Convenience functions using default interner

// Intern interns a string using the default interner
func Intern(s string) string {
    return DefaultInterner.Intern(s)
}

// InternSlice interns a slice using the default interner
func InternSlice(ss []string) []string {
    return DefaultInterner.InternSlice(ss)
}

// InternMap interns a map using the default interner
func InternMap(m map[string]string) map[string]string {
    return DefaultInterner.InternMap(m)
}
```

### 2. Common String Constants

```go
// internal/intern/common_strings.go
package intern

// Pre-intern common resource type strings
func init() {
    commonTypes := []string{
        "dashboard",
        "query",
        "control",
        "benchmark",
        "card",
        "chart",
        "container",
        "flow",
        "graph",
        "hierarchy",
        "image",
        "input",
        "node",
        "edge",
        "table",
        "text",
        "category",
        "detection",
        "detection_benchmark",
        "variable",
        "locals",
    }

    for _, t := range commonTypes {
        DefaultInterner.Intern(t)
    }
}

// Resource type constants (pre-interned)
var (
    TypeDashboard  = Intern("dashboard")
    TypeQuery      = Intern("query")
    TypeControl    = Intern("control")
    TypeBenchmark  = Intern("benchmark")
    TypeCard       = Intern("card")
    TypeChart      = Intern("chart")
    TypeContainer  = Intern("container")
    // ... etc
)
```

### 3. Integration with Index Entry

```go
// internal/resourceindex/entry.go modifications

import "github.com/turbot/powerpipe/internal/intern"

// NewIndexEntry creates an index entry with interned strings
func NewIndexEntry(resourceType, name, shortName, title, fileName string,
    startLine, endLine int, modName string) *IndexEntry {

    return &IndexEntry{
        Type:      intern.Intern(resourceType),
        Name:      intern.Intern(name),
        ShortName: intern.Intern(shortName),
        Title:     title, // Titles are usually unique, don't intern
        FileName:  intern.Intern(fileName), // File names are repeated
        StartLine: startLine,
        EndLine:   endLine,
        ModName:   intern.Intern(modName),
    }
}

// SetTags sets tags with interned keys (tag keys are often repeated)
func (e *IndexEntry) SetTags(tags map[string]string) {
    if len(tags) == 0 {
        return
    }

    e.Tags = make(map[string]string, len(tags))
    for k, v := range tags {
        // Intern keys (often repeated: service, category, etc.)
        // Values are usually unique, don't intern
        e.Tags[intern.Intern(k)] = v
    }
}
```

### 4. Integration with Scanner

```go
// internal/resourceindex/scanner.go modifications

func (s *Scanner) finalizeBlock(block *blockState, endLine int, filePath string) {
    // Use interned strings for common values
    entry := &IndexEntry{
        Type:      intern.Intern(block.blockType),
        Name:      intern.Intern(s.modName + "." + block.blockType + "." + block.name),
        ShortName: intern.Intern(block.name),
        FileName:  intern.Intern(filePath),
        StartLine: block.startLine,
        EndLine:   endLine,
        ModName:   intern.Intern(s.modName),
    }

    // Set attributes
    if title, ok := block.attributes["title"]; ok {
        entry.Title = title // Don't intern titles
    }

    // Intern tag keys
    if tags, ok := block.attributes["tags"]; ok {
        entry.Tags = parseTags(tags) // parseTags should intern keys
    }

    s.index.Add(entry)
}
```

### 5. Tests

```go
// internal/intern/string_intern_test.go
package intern

import (
    "runtime"
    "sync"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestStringInterner_Basic(t *testing.T) {
    i := NewStringInterner()

    s1 := i.Intern("hello")
    s2 := i.Intern("hello")
    s3 := i.Intern("world")

    // Same string should return same instance
    assert.True(t, &s1 == &s2 || s1 == s2)
    assert.NotEqual(t, s1, s3)

    stats := i.Stats()
    assert.Equal(t, 2, stats.UniqueStrings) // "hello" and "world"
    assert.Equal(t, int64(1), stats.Hits)   // Second "hello"
    assert.Equal(t, int64(2), stats.Misses) // First "hello" and "world"
}

func TestStringInterner_Concurrent(t *testing.T) {
    i := NewStringInterner()

    var wg sync.WaitGroup
    for j := 0; j < 100; j++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for k := 0; k < 1000; k++ {
                i.Intern("shared_string")
                i.Intern("another_string")
            }
        }(j)
    }

    wg.Wait()

    stats := i.Stats()
    assert.Equal(t, 2, stats.UniqueStrings)
    assert.Greater(t, stats.HitRate, 0.99) // Most should be hits
}

func TestStringInterner_Slice(t *testing.T) {
    i := NewStringInterner()

    input := []string{"a", "b", "a", "c", "b"}
    result := i.InternSlice(input)

    assert.Equal(t, len(input), len(result))
    assert.Equal(t, 3, i.Stats().UniqueStrings)
}

func TestStringInterner_Map(t *testing.T) {
    i := NewStringInterner()

    input := map[string]string{
        "key1": "value1",
        "key2": "value2",
    }

    // First call
    result1 := i.InternMap(input)

    // Second call with same keys
    input2 := map[string]string{
        "key1": "different",
        "key2": "values",
    }
    result2 := i.InternMap(input2)

    // Keys should be shared
    assert.Equal(t, 4, i.Stats().UniqueStrings) // key1, key2, value1, value2, different, values = 6
    _ = result1
    _ = result2
}

func TestStringInterner_MemorySavings(t *testing.T) {
    // Create many repeated strings without interning
    runtime.GC()
    var beforeNoIntern runtime.MemStats
    runtime.ReadMemStats(&beforeNoIntern)

    noInternStrings := make([]string, 10000)
    for j := 0; j < 10000; j++ {
        noInternStrings[j] = string([]byte("repeated_string_" + string(rune('A'+j%26))))
    }

    runtime.GC()
    var afterNoIntern runtime.MemStats
    runtime.ReadMemStats(&afterNoIntern)

    noInternMem := afterNoIntern.HeapAlloc - beforeNoIntern.HeapAlloc

    // Now with interning
    i := NewStringInterner()

    runtime.GC()
    var beforeIntern runtime.MemStats
    runtime.ReadMemStats(&beforeIntern)

    internStrings := make([]string, 10000)
    for j := 0; j < 10000; j++ {
        internStrings[j] = i.Intern("repeated_string_" + string(rune('A'+j%26)))
    }

    runtime.GC()
    var afterIntern runtime.MemStats
    runtime.ReadMemStats(&afterIntern)

    internMem := afterIntern.HeapAlloc - beforeIntern.HeapAlloc

    t.Logf("Without interning: %d bytes", noInternMem)
    t.Logf("With interning: %d bytes", internMem)
    t.Logf("Savings: %d bytes (%.1f%%)", noInternMem-internMem,
        float64(noInternMem-internMem)/float64(noInternMem)*100)

    // Should save memory (26 unique strings vs 10000)
    _ = noInternStrings
    _ = internStrings
}

func BenchmarkStringInterner_Intern(b *testing.B) {
    i := NewStringInterner()

    // Pre-populate with some strings
    for j := 0; j < 100; j++ {
        i.Intern(string(rune('A' + j%26)))
    }

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        i.Intern(string(rune('A' + n%26)))
    }
}

func BenchmarkStringInterner_Concurrent(b *testing.B) {
    i := NewStringInterner()

    b.RunParallel(func(pb *testing.PB) {
        n := 0
        for pb.Next() {
            i.Intern(string(rune('A' + n%26)))
            n++
        }
    })
}
```

## Acceptance Criteria

- [ ] StringInterner correctly deduplicates strings
- [ ] Thread-safe for concurrent access
- [ ] InternSlice and InternMap work correctly
- [ ] Integration with index entry creation
- [ ] Integration with scanner
- [ ] Common type strings pre-interned
- [ ] Hit rate tracking for monitoring
- [ ] Memory savings measurable
- [ ] Performance benchmarks acceptable (< 100ns per intern)
- [ ] All tests pass

## Notes

- String interning has diminishing returns for unique strings
- Focus on repeated values: types, file names, mod names, tag keys
- Don't intern titles, descriptions, SQL (usually unique)
- Consider weak references to allow GC of unused strings
- May want to limit intern table size to prevent unbounded growth
- This is a lower priority optimization - implement if time permits
