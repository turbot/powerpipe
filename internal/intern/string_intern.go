package intern

import (
	"sync"
	"sync/atomic"
)

// StringInterner deduplicates strings to reduce memory usage.
// Thread-safe for concurrent access with read-optimized locking.
type StringInterner struct {
	mu      sync.RWMutex
	strings map[string]string
	hits    atomic.Int64
	misses  atomic.Int64
}

// NewStringInterner creates a new interner.
func NewStringInterner() *StringInterner {
	return &StringInterner{
		strings: make(map[string]string),
	}
}

// DefaultInterner is the global string interner for common use.
var DefaultInterner = NewStringInterner()

// Intern returns a canonical version of the string.
// If the string was seen before, returns the original instance.
// This reduces memory by sharing string backing arrays.
func (i *StringInterner) Intern(s string) string {
	if s == "" {
		return ""
	}

	// Fast path: check if already interned (read lock only)
	i.mu.RLock()
	if interned, ok := i.strings[s]; ok {
		i.hits.Add(1)
		i.mu.RUnlock()
		return interned
	}
	i.mu.RUnlock()

	// Slow path: add to intern table (write lock)
	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after acquiring write lock
	if interned, ok := i.strings[s]; ok {
		i.hits.Add(1)
		return interned
	}

	// Intern the string - make a copy to ensure we own the backing array
	// This is important when strings are substrings of larger buffers
	interned := string([]byte(s))
	i.strings[interned] = interned
	i.misses.Add(1)
	return interned
}

// InternSlice interns all strings in a slice and returns a new slice.
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

// InternMapKeys interns all keys in a map (values are left as-is).
// Use this for maps with repeated keys (like tags).
func (i *StringInterner) InternMapKeys(m map[string]string) map[string]string {
	if len(m) == 0 {
		return m
	}

	result := make(map[string]string, len(m))
	for k, v := range m {
		result[i.Intern(k)] = v
	}
	return result
}

// InternMap interns all keys and values in a map.
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

// Stats returns interning statistics.
func (i *StringInterner) Stats() InternStats {
	i.mu.RLock()
	defer i.mu.RUnlock()

	hits := i.hits.Load()
	misses := i.misses.Load()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	// Estimate memory saved: each hit saves the string length
	var savedBytes int64
	if len(i.strings) > 0 && hits > 0 {
		var totalLen int64
		for s := range i.strings {
			totalLen += int64(len(s))
		}
		avgLen := totalLen / int64(len(i.strings))
		savedBytes = avgLen * hits
	}

	return InternStats{
		UniqueStrings: len(i.strings),
		Hits:          hits,
		Misses:        misses,
		HitRate:       hitRate,
		SavedBytes:    savedBytes,
	}
}

// Size returns the number of unique strings interned.
func (i *StringInterner) Size() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.strings)
}

// Clear clears the intern table and resets statistics.
func (i *StringInterner) Clear() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.strings = make(map[string]string)
	i.hits.Store(0)
	i.misses.Store(0)
}

// InternStats contains statistics about interning operations.
type InternStats struct {
	UniqueStrings int     // Number of unique strings stored
	Hits          int64   // Number of cache hits (string already interned)
	Misses        int64   // Number of cache misses (new string added)
	HitRate       float64 // Ratio of hits to total requests
	SavedBytes    int64   // Estimated bytes saved by interning
}

// Convenience functions using the default interner.

// Intern interns a string using the default interner.
func Intern(s string) string {
	return DefaultInterner.Intern(s)
}

// InternSlice interns a slice using the default interner.
func InternSlice(ss []string) []string {
	return DefaultInterner.InternSlice(ss)
}

// InternMapKeys interns map keys using the default interner.
func InternMapKeys(m map[string]string) map[string]string {
	return DefaultInterner.InternMapKeys(m)
}

// InternMap interns a map using the default interner.
func InternMap(m map[string]string) map[string]string {
	return DefaultInterner.InternMap(m)
}

// Stats returns statistics for the default interner.
func Stats() InternStats {
	return DefaultInterner.Stats()
}

// Clear clears the default interner.
func Clear() {
	DefaultInterner.Clear()
}

// Reset resets the default interner to a fresh state.
// This also pre-interns common strings.
func Reset() {
	DefaultInterner.Clear()
	preInternCommonStrings(DefaultInterner)
}
