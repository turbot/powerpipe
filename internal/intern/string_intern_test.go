package intern

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringInterner_Basic(t *testing.T) {
	i := NewStringInterner()

	s1 := i.Intern("hello")
	s2 := i.Intern("hello")
	s3 := i.Intern("world")

	// Same string should return same value
	assert.Equal(t, s1, s2)
	assert.Equal(t, "hello", s1)
	assert.NotEqual(t, s1, s3)

	stats := i.Stats()
	assert.Equal(t, 2, stats.UniqueStrings) // "hello" and "world"
	assert.Equal(t, int64(1), stats.Hits)   // Second "hello"
	assert.Equal(t, int64(2), stats.Misses) // First "hello" and "world"
}

func TestStringInterner_EmptyString(t *testing.T) {
	i := NewStringInterner()

	s := i.Intern("")
	assert.Equal(t, "", s)

	// Empty strings shouldn't be added to the intern table
	stats := i.Stats()
	assert.Equal(t, 0, stats.UniqueStrings)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
}

func TestStringInterner_Concurrent(t *testing.T) {
	i := NewStringInterner()

	var wg sync.WaitGroup
	goroutines := 100
	iterations := 1000

	for j := 0; j < goroutines; j++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for k := 0; k < iterations; k++ {
				i.Intern("shared_string")
				i.Intern("another_string")
				// Also intern some unique strings per goroutine
				if k%100 == 0 {
					i.Intern(string(rune('A' + id%26)))
				}
			}
		}(j)
	}

	wg.Wait()

	stats := i.Stats()
	// Should have 2 shared + up to 26 unique per-goroutine strings
	assert.GreaterOrEqual(t, stats.UniqueStrings, 2)
	assert.LessOrEqual(t, stats.UniqueStrings, 2+26)
	// Most operations should be hits (except first occurrence of each)
	assert.Greater(t, stats.HitRate, 0.9)
}

func TestStringInterner_InternSlice(t *testing.T) {
	i := NewStringInterner()

	input := []string{"a", "b", "a", "c", "b"}
	result := i.InternSlice(input)

	assert.Equal(t, len(input), len(result))
	assert.Equal(t, input, result)
	assert.Equal(t, 3, i.Stats().UniqueStrings) // a, b, c
}

func TestStringInterner_InternSliceEmpty(t *testing.T) {
	i := NewStringInterner()

	var empty []string
	result := i.InternSlice(empty)
	assert.Empty(t, result)
	assert.Equal(t, 0, i.Stats().UniqueStrings)
}

func TestStringInterner_InternMapKeys(t *testing.T) {
	i := NewStringInterner()

	input := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	result := i.InternMapKeys(input)

	// Values should be preserved as-is
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])

	// Only keys should be interned
	assert.Equal(t, 2, i.Stats().UniqueStrings) // key1, key2

	// Interning same keys again should hit cache
	input2 := map[string]string{
		"key1": "different_value",
		"key2": "another_value",
	}
	i.InternMapKeys(input2)

	stats := i.Stats()
	assert.Equal(t, 2, stats.UniqueStrings) // Still just 2
	assert.Equal(t, int64(2), stats.Hits)   // key1 and key2 hit
}

func TestStringInterner_InternMap(t *testing.T) {
	i := NewStringInterner()

	input := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	result := i.InternMap(input)

	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, 4, i.Stats().UniqueStrings) // key1, key2, value1, value2
}

func TestStringInterner_InternMapEmpty(t *testing.T) {
	i := NewStringInterner()

	var empty map[string]string
	result := i.InternMapKeys(empty)
	assert.Empty(t, result)

	result = i.InternMap(empty)
	assert.Empty(t, result)

	assert.Equal(t, 0, i.Stats().UniqueStrings)
}

func TestStringInterner_Clear(t *testing.T) {
	i := NewStringInterner()

	i.Intern("hello")
	i.Intern("world")
	assert.Equal(t, 2, i.Stats().UniqueStrings)

	i.Clear()

	stats := i.Stats()
	assert.Equal(t, 0, stats.UniqueStrings)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
}

func TestStringInterner_Size(t *testing.T) {
	i := NewStringInterner()

	assert.Equal(t, 0, i.Size())

	i.Intern("a")
	assert.Equal(t, 1, i.Size())

	i.Intern("b")
	assert.Equal(t, 2, i.Size())

	i.Intern("a") // duplicate
	assert.Equal(t, 2, i.Size())
}

func TestStringInterner_HitRate(t *testing.T) {
	i := NewStringInterner()

	// First intern - miss
	i.Intern("test")
	assert.Equal(t, float64(0), i.Stats().HitRate)

	// Second intern - hit
	i.Intern("test")
	assert.Equal(t, 0.5, i.Stats().HitRate)

	// Third intern - hit
	i.Intern("test")
	assert.InDelta(t, 0.667, i.Stats().HitRate, 0.01)
}

func TestDefaultInterner(t *testing.T) {
	// Reset to clear any state from other tests
	Clear()

	s1 := Intern("default_test")
	s2 := Intern("default_test")
	assert.Equal(t, s1, s2)

	stats := Stats()
	assert.GreaterOrEqual(t, stats.UniqueStrings, 1)
}

func TestCommonStrings_PreInterned(t *testing.T) {
	// Common type strings should be pre-interned
	assert.Equal(t, "dashboard", TypeDashboard)
	assert.Equal(t, "benchmark", TypeBenchmark)
	assert.Equal(t, "control", TypeControl)
	assert.Equal(t, "query", TypeQuery)
	assert.Equal(t, "card", TypeCard)
	assert.Equal(t, "chart", TypeChart)
	assert.Equal(t, "container", TypeContainer)
	assert.Equal(t, "flow", TypeFlow)
	assert.Equal(t, "graph", TypeGraph)
	assert.Equal(t, "hierarchy", TypeHierarchy)
	assert.Equal(t, "image", TypeImage)
	assert.Equal(t, "input", TypeInput)
	assert.Equal(t, "node", TypeNode)
	assert.Equal(t, "edge", TypeEdge)
	assert.Equal(t, "table", TypeTable)
	assert.Equal(t, "text", TypeText)
	assert.Equal(t, "category", TypeCategory)
	assert.Equal(t, "detection", TypeDetection)
	assert.Equal(t, "detection_benchmark", TypeDetectionBenchmark)
	assert.Equal(t, "variable", TypeVariable)
	assert.Equal(t, "with", TypeWith)
}

func TestCommonStrings_InternReusesPreInterned(t *testing.T) {
	// Create a fresh interner
	i := NewStringInterner()
	preInternCommonStrings(i)

	hitsBefore := i.Stats().Hits

	// Interning "dashboard" should hit cache
	result := i.Intern("dashboard")
	assert.Equal(t, "dashboard", result)

	hitsAfter := i.Stats().Hits
	assert.Equal(t, hitsBefore+1, hitsAfter)
}

func TestStringInterner_SubstringHandling(t *testing.T) {
	i := NewStringInterner()

	// When strings are substrings of larger buffers, we should
	// make a copy to avoid holding onto large buffers
	largeString := "this is a very long string that we don't want to keep"
	substring := largeString[0:4] // "this"

	interned := i.Intern(substring)
	assert.Equal(t, "this", interned)

	// Clear reference to large string
	largeString = ""
	runtime.GC()

	// Interned string should still work
	assert.Equal(t, "this", interned)
}

func TestReset(t *testing.T) {
	// Add some strings
	Intern("custom_string_1")
	Intern("custom_string_2")

	// Reset should clear and re-initialize with common strings
	Reset()

	// Common strings should still work (and be hits)
	s := Intern("dashboard")
	assert.Equal(t, "dashboard", s)

	// Stats show we have at least the common strings
	stats := Stats()
	assert.GreaterOrEqual(t, stats.UniqueStrings, 20) // Pre-interned count
}

func TestStringInterner_MemorySavings(t *testing.T) {
	// Create strings without interning
	var noInternStrings []string
	for j := 0; j < 1000; j++ {
		// Create 26 different strings, repeated many times
		s := "repeated_string_" + string(rune('A'+j%26))
		noInternStrings = append(noInternStrings, string([]byte(s))) // Force copy
	}

	// Now with interning
	i := NewStringInterner()
	var internStrings []string
	for j := 0; j < 1000; j++ {
		s := "repeated_string_" + string(rune('A'+j%26))
		internStrings = append(internStrings, i.Intern(s))
	}

	stats := i.Stats()

	t.Logf("Without interning: 1000 strings")
	t.Logf("With interning: %d unique strings", stats.UniqueStrings)
	t.Logf("Hit rate: %.1f%%", stats.HitRate*100)
	t.Logf("Estimated savings: %d bytes", stats.SavedBytes)

	// Should have only 26 unique strings
	assert.Equal(t, 26, stats.UniqueStrings)
	// Should have high hit rate
	assert.Greater(t, stats.HitRate, 0.95)

	// Keep variables alive for proper measurement
	require.NotEmpty(t, noInternStrings)
	require.NotEmpty(t, internStrings)
}

// Benchmarks

func BenchmarkStringInterner_Intern_Hit(b *testing.B) {
	i := NewStringInterner()

	// Pre-populate
	i.Intern("benchmark_string")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i.Intern("benchmark_string")
	}
}

func BenchmarkStringInterner_Intern_Miss(b *testing.B) {
	i := NewStringInterner()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Create unique strings to always miss
		i.Intern(string(rune(n % 65536)))
	}
}

func BenchmarkStringInterner_Concurrent(b *testing.B) {
	i := NewStringInterner()

	// Pre-populate with some strings
	for j := 0; j < 100; j++ {
		i.Intern(string(rune('A' + j%26)))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		n := 0
		for pb.Next() {
			i.Intern(string(rune('A' + n%26)))
			n++
		}
	})
}

func BenchmarkInternSlice(b *testing.B) {
	i := NewStringInterner()
	slice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	// Pre-intern
	i.InternSlice(slice)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i.InternSlice(slice)
	}
}

func BenchmarkInternMapKeys(b *testing.B) {
	i := NewStringInterner()
	m := map[string]string{
		"service":  "aws",
		"category": "security",
		"type":     "benchmark",
		"plugin":   "terraform",
	}

	// Pre-intern keys
	i.InternMapKeys(m)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i.InternMapKeys(m)
	}
}
