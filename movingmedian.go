// Package movingmedian computes the median of a windowed stream of data.
package movingmedian

import (
	"container/heap"
	"fmt"
)

type item struct {
	f         float64
	heapIndex int
}

type itemHeap []*item

func (h itemHeap) Len() int { return len(h) }
func (h itemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *itemHeap) Push(x interface{}) {
	e, ok := x.(*item)
	if !ok {
		panic(fmt.Sprintf("%T is not an *item", x)) //nolint:forbidigo
	}
	e.heapIndex = len(*h)
	*h = append(*h, e)
}

func (h *itemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type minItemHeap struct {
	itemHeap
}

func (h minItemHeap) Less(i, j int) bool { return h.itemHeap[i].f < h.itemHeap[j].f }

type maxItemHeap struct {
	itemHeap
}

func (h maxItemHeap) Less(i, j int) bool { return h.itemHeap[i].f > h.itemHeap[j].f }

// MovingMedian computes the moving median of a windowed stream of numbers.
type MovingMedian struct {
	queueIndex int
	nitems     int
	queue      []item
	maxHeap    maxItemHeap
	minHeap    minItemHeap
}

// NewMovingMedian returns a MovingMedian with the given window size.
func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		queue:   make([]item, size),
		maxHeap: maxItemHeap{make([]*item, 0, size/2+1)}, // Pre-allocate with capacity
		minHeap: minItemHeap{make([]*item, 0, size/2+1)},
	}
	heap.Init(&m.maxHeap)
	heap.Init(&m.minHeap)
	return m
}

// Push adds an element to the stream, removing old data which has expired from the window.  It runs in O(log windowSize).
func (m *MovingMedian) Push(v float64) {
	// Special case optimization for size 1
	if len(m.queue) == 1 {
		m.queue[0].f = v
		return
	}

	// Use pointer arithmetic for queue management
	itemPtr := &m.queue[m.queueIndex]
	m.queueIndex++
	if m.queueIndex >= len(m.queue) {
		m.queueIndex = 0
	}

	// Update value
	itemPtr.f = v

	// Fast path for initial window filling
	if m.nitems < len(m.queue) {
		m.nitems++

		// Simple balancing for initial fill
		if m.minHeap.Len() == 0 || v > m.minHeap.itemHeap[0].f {
			heap.Push(&m.minHeap, itemPtr)
			if m.minHeap.Len() > m.maxHeap.Len()+1 {
				moveItem := heap.Pop(&m.minHeap)
				heap.Push(&m.maxHeap, moveItem)
			}
			return
		}

		heap.Push(&m.maxHeap, itemPtr)
		if m.maxHeap.Len() > m.minHeap.Len()+1 {
			moveItem := heap.Pop(&m.maxHeap)
			heap.Push(&m.minHeap, moveItem)
		}
		return
	}

	// Main path for full window updates
	minAbove := m.minHeap.itemHeap[0].f
	maxBelow := m.maxHeap.itemHeap[0].f

	// Check if item is in min heap
	if itemPtr.heapIndex < m.minHeap.Len() && itemPtr == m.minHeap.itemHeap[itemPtr.heapIndex] {
		if v >= maxBelow {
			heap.Fix(&m.minHeap, itemPtr.heapIndex)
			return
		}
		rotate(&m.maxHeap, &m.minHeap, m.maxHeap.itemHeap, m.minHeap.itemHeap, itemPtr)
		return
	}

	// Item must be in max heap
	if v <= minAbove {
		heap.Fix(&m.maxHeap, itemPtr.heapIndex)
		return
	}
	rotate(&m.minHeap, &m.maxHeap, m.minHeap.itemHeap, m.maxHeap.itemHeap, itemPtr)
}

func rotate(heapA, heapB heap.Interface, itemHeapA, itemHeapB itemHeap, itemPtr *item) {
	moveItem := itemHeapA[0]
	moveItem.heapIndex = itemPtr.heapIndex
	itemHeapB[itemPtr.heapIndex] = moveItem
	itemHeapA[0] = itemPtr
	heap.Fix(heapB, itemPtr.heapIndex)
	itemPtr.heapIndex = 0
	heap.Fix(heapA, 0)
}

// Median returns the current value of the median from the window.
func (m *MovingMedian) Median() float64 {
	if len(m.queue) == 1 {
		return m.queue[0].f
	}

	minLen := m.minHeap.Len()
	maxLen := m.maxHeap.Len()

	if maxLen == minLen {
		if maxLen == 0 {
			return 0.0
		}
		return (m.maxHeap.itemHeap[0].f + m.minHeap.itemHeap[0].f) / 2
	}

	if maxLen > minLen {
		return m.maxHeap.itemHeap[0].f
	}

	return m.minHeap.itemHeap[0].f
}
