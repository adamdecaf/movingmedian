package movingmedian

import (
	"container/heap"
	"math"
)

type item struct {
	f   float64
	idx int
}

type itemHeap []*item

func (h itemHeap) Len() int { return len(h) }
func (h itemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *itemHeap) Push(x interface{}) {
	e := x.(*item)
	e.idx = len(*h)
	*h = append(*h, e)
}

func (h *itemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type minFloat64Heap struct {
	itemHeap
}

func (h minFloat64Heap) Less(i, j int) bool { return h.itemHeap[i].f < h.itemHeap[j].f }

type maxFloat64Heap struct {
	itemHeap
}

func (h maxFloat64Heap) Less(i, j int) bool { return h.itemHeap[i].f > h.itemHeap[j].f }

type MovingMedian struct {
	queueIndex int
	nitems     int
	queue      []item
	maxHeap    maxFloat64Heap
	minHeap    minFloat64Heap
}

func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		queue:   make([]item, size),
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

	if m.nitems == len(m.queue) {
		old := &m.queue[m.queueIndex]
		heapIndex := old.idx

		if heapIndex < m.minHeap.Len() && old == m.minHeap.itemHeap[heapIndex] {
			heap.Remove(&m.minHeap, heapIndex)
		} else {
			heap.Remove(&m.maxHeap, heapIndex)
		}
	} else {
		m.nitems++
	}

	m.queue[m.queueIndex] = item{f: v}
	e := &m.queue[m.queueIndex]

	m.queueIndex++

	if m.queueIndex >= len(m.queue) {
		m.queueIndex = 0
	}

	if m.minHeap.Len() == 0 ||
		v > m.minHeap.itemHeap[0].f {
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

	if (m.nitems % 2) == 0 {
		return (m.maxHeap.itemHeap[0].f + m.minHeap.itemHeap[0].f) / 2
	}

	if m.maxHeap.Len() > m.minHeap.Len() {
		return m.maxHeap.itemHeap[0].f
	}

	return m.minHeap.itemHeap[0].f
}
