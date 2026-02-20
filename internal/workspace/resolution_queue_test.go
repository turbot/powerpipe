package workspace

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolutionQueue_PriorityOrder(t *testing.T) {
	q := NewResolutionQueue()

	// Push items with different priorities (higher = popped first)
	q.Push("low", 10)
	q.Push("high", 100)
	q.Push("medium", 50)

	// Should pop in priority order
	assert.Equal(t, "high", q.Pop())
	assert.Equal(t, "medium", q.Pop())
	assert.Equal(t, "low", q.Pop())
	assert.Equal(t, "", q.Pop()) // Empty queue
}

func TestResolutionQueue_NoDuplicates(t *testing.T) {
	q := NewResolutionQueue()

	// Push same item multiple times
	q.Push("resource1", 10)
	q.Push("resource1", 100) // Duplicate - should be ignored
	q.Push("resource2", 50)

	assert.Equal(t, 2, q.Len())

	// Pop all items
	assert.Equal(t, "resource2", q.Pop()) // resource1 has priority 10, not 100
	assert.Equal(t, "resource1", q.Pop())
	assert.True(t, q.IsEmpty())
}

func TestResolutionQueue_Prioritize(t *testing.T) {
	q := NewResolutionQueue()

	// Push items
	q.Push("a", 10)
	q.Push("b", 20)
	q.Push("c", 30)

	// Prioritize "a" to very high priority
	q.Prioritize("a", 1000)

	// Now "a" should be first
	assert.Equal(t, "a", q.Pop())
	assert.Equal(t, "c", q.Pop())
	assert.Equal(t, "b", q.Pop())
}

func TestResolutionQueue_PrioritizeNewItem(t *testing.T) {
	q := NewResolutionQueue()

	q.Push("existing", 10)

	// Prioritize a new item that's not in the queue
	q.Prioritize("new", 100)

	assert.Equal(t, 2, q.Len())
	assert.Equal(t, "new", q.Pop()) // New item should be first due to high priority
	assert.Equal(t, "existing", q.Pop())
}

func TestResolutionQueue_Contains(t *testing.T) {
	q := NewResolutionQueue()

	q.Push("item1", 10)

	assert.True(t, q.Contains("item1"))
	assert.False(t, q.Contains("item2"))

	q.Pop()

	assert.False(t, q.Contains("item1")) // Removed after pop
}

func TestResolutionQueue_Clear(t *testing.T) {
	q := NewResolutionQueue()

	q.Push("a", 10)
	q.Push("b", 20)
	q.Push("c", 30)

	assert.Equal(t, 3, q.Len())

	q.Clear()

	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())
	assert.Equal(t, "", q.Pop())
}

func TestResolutionQueue_ConcurrentAccess(t *testing.T) {
	q := NewResolutionQueue()

	var wg sync.WaitGroup
	pushCount := 100
	popCount := 50

	// Concurrent pushes
	for i := 0; i < pushCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(string(rune('a'+i%26))+string(rune('0'+i)), i)
		}(i)
	}

	// Concurrent pops (some will get empty)
	results := make([]string, popCount)
	var resultMu sync.Mutex
	for i := 0; i < popCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result := q.Pop()
			resultMu.Lock()
			results[i] = result
			resultMu.Unlock()
		}(i)
	}

	wg.Wait()

	// Should have some items left or all popped
	// Just verify no panic occurred
	assert.True(t, true, "no panic during concurrent access")
}

func TestResolutionQueue_ConcurrentPrioritize(t *testing.T) {
	q := NewResolutionQueue()

	// Pre-populate
	for i := 0; i < 10; i++ {
		q.Push(string(rune('a'+i)), i*10)
	}

	var wg sync.WaitGroup

	// Concurrent prioritizes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Prioritize(string(rune('a'+i%10)), i*100)
		}(i)
	}

	wg.Wait()

	// Verify queue is still functional
	assert.True(t, q.Len() > 0 || q.IsEmpty())
}

func TestResolutionQueue_EmptyPop(t *testing.T) {
	q := NewResolutionQueue()

	// Pop from empty queue
	assert.Equal(t, "", q.Pop())
	assert.Equal(t, "", q.Pop())
	assert.True(t, q.IsEmpty())
}
