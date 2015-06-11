package movingmedian

import (
	"container/heap"
	"math"
)

type elt struct {
	f   float64
	idx int
}

type float64Heap []*elt

func (h float64Heap) Len() int { return len(h) }
func (h float64Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *float64Heap) Push(x interface{}) {
	e := x.(*elt)
	e.idx = len(*h)
	*h = append(*h, e)
}

func (h *float64Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type minFloat64Heap struct {
	float64Heap
}

func (h minFloat64Heap) Less(i, j int) bool { return h.float64Heap[i].f < h.float64Heap[j].f }

type maxFloat64Heap struct {
	float64Heap
}

func (h maxFloat64Heap) Less(i, j int) bool { return h.float64Heap[i].f > h.float64Heap[j].f }

type MovingMedian struct {
	idx     int
	nelts   int
	queue   []elt
	maxHeap maxFloat64Heap
	minHeap minFloat64Heap
}

func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		queue:   make([]elt, size),
		maxHeap: maxFloat64Heap{},
		minHeap: minFloat64Heap{},
	}

	heap.Init(&m.maxHeap)
	heap.Init(&m.minHeap)
	return m
}

func (m *MovingMedian) balanceHeaps() {
	if m.maxHeap.Len() > (m.minHeap.Len() + 1) {
		moveItem := heap.Pop(&m.maxHeap)
		heap.Push(&m.minHeap, moveItem)
	} else if m.minHeap.Len() > (m.maxHeap.Len() + 1) {
		moveItem := heap.Pop(&m.minHeap)
		heap.Push(&m.maxHeap, moveItem)
	}

}

func (m *MovingMedian) Push(v float64) {

	if m.nelts >= len(m.queue) {
		old := &m.queue[m.idx]

		if old.idx < m.minHeap.Len() && old == m.minHeap.float64Heap[old.idx] {
			heap.Remove(&m.minHeap, old.idx)
		} else {
			heap.Remove(&m.maxHeap, old.idx)
		}
	}

	m.queue[m.idx] = elt{f: v}
	e := &m.queue[m.idx]

	m.nelts++
	m.idx++

	if m.idx >= len(m.queue) {
		m.idx = 0
	}

	if m.minHeap.Len() == 0 ||
		v > m.minHeap.float64Heap[0].f {
		heap.Push(&m.minHeap, e)
	} else {
		heap.Push(&m.maxHeap, e)
	}

	m.balanceHeaps()
}

func (m *MovingMedian) Median() float64 {
	if len(m.queue) == 0 {
		return math.NaN()
	}

	wsize := m.nelts
	if m.nelts > len(m.queue) {
		wsize = len(m.queue)
	}

	if (wsize % 2) == 0 {
		return (m.maxHeap.float64Heap[0].f + m.minHeap.float64Heap[0].f) / 2
	}

	if m.maxHeap.Len() > m.minHeap.Len() {
		return m.maxHeap.float64Heap[0].f
	}

	return m.minHeap.float64Heap[0].f
}
