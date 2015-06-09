package movingmedian

import (
	"container/heap"
	"math"
)

type HeapInterface interface {
	heap.Interface
	RemoveValue(float64) bool
}

type Heap []float64

func (h Heap) Len() int { return len(h) }
func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MinHeap struct {
	Heap
}

func (h MinHeap) Less(i, j int) bool { return h.Heap[i] < h.Heap[j] }

func (h *MinHeap) RemoveValue(x float64) bool {
	for i, v := range h.Heap {
		if v == x {
			heap.Remove(h, i)
			return true
		}
	}

	return false
}

type MaxHeap struct {
	Heap
}

func (h MaxHeap) Less(i, j int) bool { return h.Heap[i] > h.Heap[j] }

func (h *MaxHeap) RemoveValue(x float64) bool {
	for i, v := range h.Heap {
		if v == x {
			heap.Remove(h, i)
			return true
		}
	}

	return false
}

type MovingMedian struct {
	size    int
	queue   []float64
	maxHeap MaxHeap
	minHeap MinHeap
}

func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		size:    size,
		queue:   make([]float64, 0),
		maxHeap: MaxHeap{},
		minHeap: MinHeap{},
	}

	heap.Init(&m.maxHeap)
	heap.Init(&m.minHeap)
	return m
}

func (m *MovingMedian) Push(v float64) {
	m.queue = append(m.queue, v)

	if m.minHeap.Len() == 0 ||
		v > m.minHeap.Heap[0] {
		heap.Push(&m.minHeap, v)
	} else {
		heap.Push(&m.maxHeap, v)
	}

	if len(m.queue) > m.size {
		outItem := m.queue[0]
		m.queue = m.queue[1:len(m.queue)]
		if !m.minHeap.RemoveValue(outItem) {
			m.maxHeap.RemoveValue(outItem)
		}
	}

	if m.maxHeap.Len() > (m.minHeap.Len() + 1) {
		moveItem := heap.Pop(&m.maxHeap).(float64)
		heap.Push(&m.minHeap, moveItem)
	} else if m.minHeap.Len() > (m.maxHeap.Len() + 1) {
		moveItem := heap.Pop(&m.minHeap).(float64)
		heap.Push(&m.maxHeap, moveItem)
	}
}

func (m MovingMedian) Median() float64 {
	if len(m.queue) == 0 {
		return math.NaN()
	}

	if (len(m.queue) % 2) == 0 {
		return (m.maxHeap.Heap[0] + m.minHeap.Heap[0]) / 2
	}

	if m.maxHeap.Len() > m.minHeap.Len() {
		return m.maxHeap.Heap[0]
	}

	return m.minHeap.Heap[0]
}
