package workspace

import (
	"container/heap"
	"sync"
)

// ResolutionQueue is a thread-safe priority queue for resolution requests.
// Higher priority items are processed first.
type ResolutionQueue struct {
	items priorityQueue
	seen  map[string]bool // Track items already in queue
	mu    sync.Mutex
}

// queueItem represents an item in the priority queue.
type queueItem struct {
	name     string
	priority int
	index    int // Heap index (maintained by heap.Interface)
}

// priorityQueue implements heap.Interface for queueItems.
type priorityQueue []*queueItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// Higher priority first
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*queueItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// NewResolutionQueue creates a new priority queue.
func NewResolutionQueue() *ResolutionQueue {
	q := &ResolutionQueue{
		items: make(priorityQueue, 0),
		seen:  make(map[string]bool),
	}
	heap.Init(&q.items)
	return q
}

// Push adds a resource to the queue with the given priority.
// Does nothing if the resource is already in the queue.
func (q *ResolutionQueue) Push(name string, priority int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Skip if already in queue
	if q.seen[name] {
		return
	}

	q.seen[name] = true
	heap.Push(&q.items, &queueItem{name: name, priority: priority})
}

// Pop removes and returns the highest priority resource name.
// Returns empty string if queue is empty.
func (q *ResolutionQueue) Pop() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return ""
	}

	item := heap.Pop(&q.items).(*queueItem)
	delete(q.seen, item.name)
	return item.name
}

// Prioritize moves a resource to higher priority if it exists in the queue,
// or adds it with high priority if not present.
func (q *ResolutionQueue) Prioritize(name string, newPriority int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// If not in queue, add with high priority
	if !q.seen[name] {
		q.seen[name] = true
		heap.Push(&q.items, &queueItem{name: name, priority: newPriority})
		return
	}

	// Find and update priority
	for _, item := range q.items {
		if item.name == name {
			item.priority = newPriority
			heap.Fix(&q.items, item.index)
			return
		}
	}
}

// IsEmpty returns true if the queue has no items.
func (q *ResolutionQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}

// Len returns the number of items in the queue.
func (q *ResolutionQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Contains returns true if the resource is in the queue.
func (q *ResolutionQueue) Contains(name string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.seen[name]
}

// Clear removes all items from the queue.
func (q *ResolutionQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make(priorityQueue, 0)
	q.seen = make(map[string]bool)
	heap.Init(&q.items)
}
