package movingmedian

import (
	"container/heap"
	"math"
)

type float64HeapInterface interface {
	heap.Interface
	Data() []float64
}

type float64Heap []float64

func (h float64Heap) Len() int { return len(h) }
func (h float64Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *float64Heap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

func (h *float64Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h float64Heap) Data() []float64 {
	return []float64(h)
}

func removeFromHeap(h float64HeapInterface, x float64) bool {
	for i, v := range h.Data() {
		if v == x {
			heap.Remove(h, i)
			return true
		}
	}

	return false
}

type minFloat64Heap struct {
	float64Heap
}

func (h minFloat64Heap) Less(i, j int) bool { return h.float64Heap[i] < h.float64Heap[j] }

type maxFloat64Heap struct {
	float64Heap
}

func (h maxFloat64Heap) Less(i, j int) bool { return h.float64Heap[i] > h.float64Heap[j] }

type MovingMedian struct {
	size    int
	queue   []float64
	maxHeap maxFloat64Heap
	minHeap minFloat64Heap
}

func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		size:    size,
		queue:   make([]float64, 0),
		maxHeap: maxFloat64Heap{},
		minHeap: minFloat64Heap{},
	}

	heap.Init(&m.maxHeap)
	heap.Init(&m.minHeap)
	return m
}

func (m *MovingMedian) Push(v float64) {
	m.queue = append(m.queue, v)

	if m.minHeap.Len() == 0 ||
		v > m.minHeap.float64Heap[0] {
		heap.Push(&m.minHeap, v)
	} else {
		heap.Push(&m.maxHeap, v)
	}

	if len(m.queue) > m.size {
		outItem := m.queue[0]
		m.queue = m.queue[1:len(m.queue)]
		if outItem >= m.minHeap.float64Heap[0] {
			removeFromHeap(&m.minHeap, outItem)
		} else {
			removeFromHeap(&m.maxHeap, outItem)
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
		return (m.maxHeap.float64Heap[0] + m.minHeap.float64Heap[0]) / 2
	}

	if m.maxHeap.Len() > m.minHeap.Len() {
		return m.maxHeap.float64Heap[0]
	}

	return m.minHeap.float64Heap[0]
}
