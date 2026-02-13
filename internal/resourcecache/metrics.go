package resourcecache

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MetricsCollector tracks cache performance over time
type MetricsCollector struct {
	mu         sync.Mutex
	cache      *Cache
	interval   time.Duration
	samples    []CacheStats
	maxSamples int
	stopCh     chan struct{}
	running    bool
}

// NewMetricsCollector creates a new metrics collector for the given cache
func NewMetricsCollector(cache *Cache, interval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		cache:      cache,
		interval:   interval,
		maxSamples: 1000,
		stopCh:     make(chan struct{}),
	}
}

// Start begins periodic metrics collection
func (m *MetricsCollector) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.Collect()
			case <-m.stopCh:
				return
			}
		}
	}()
}

// Stop stops periodic metrics collection
func (m *MetricsCollector) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		close(m.stopCh)
		m.running = false
		m.stopCh = make(chan struct{}) // Reset for potential restart
	}
}

// Collect takes a snapshot of current cache stats
func (m *MetricsCollector) Collect() {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := m.cache.Stats()
	m.samples = append(m.samples, stats)

	// Keep bounded
	if len(m.samples) > m.maxSamples {
		m.samples = m.samples[len(m.samples)-m.maxSamples:]
	}
}

// Samples returns all collected samples
func (m *MetricsCollector) Samples() []CacheStats {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]CacheStats, len(m.samples))
	copy(result, m.samples)
	return result
}

// Report returns a human-readable report of cache statistics
func (m *MetricsCollector) Report() string {
	stats := m.cache.Stats()

	var b strings.Builder
	b.WriteString("=== Resource Cache Metrics ===\n")
	b.WriteString(fmt.Sprintf("Entries:      %d\n", stats.Entries))
	b.WriteString(fmt.Sprintf("Memory:       %s / %s (%.1f%%)\n",
		formatBytes(stats.MemoryBytes),
		formatBytes(stats.MaxMemory),
		float64(stats.MemoryBytes)/float64(stats.MaxMemory)*100))
	b.WriteString(fmt.Sprintf("Hit Rate:     %.1f%%\n", stats.HitRate*100))
	b.WriteString(fmt.Sprintf("Hits:         %d\n", stats.Hits))
	b.WriteString(fmt.Sprintf("Misses:       %d\n", stats.Misses))
	b.WriteString(fmt.Sprintf("Evictions:    %d\n", stats.Evictions))

	return b.String()
}

// JSON returns the current stats as JSON
func (m *MetricsCollector) JSON() ([]byte, error) {
	return json.Marshal(m.cache.Stats())
}

// formatBytes formats a byte count as a human-readable string
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// ResetStats resets the hit/miss/eviction counters (for testing or periodic reset)
func (m *MetricsCollector) ResetStats() {
	m.cache.mu.Lock()
	defer m.cache.mu.Unlock()

	m.cache.hits = 0
	m.cache.misses = 0
	m.cache.evictions = 0

	m.mu.Lock()
	m.samples = nil
	m.mu.Unlock()
}
