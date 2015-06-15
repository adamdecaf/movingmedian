package movingmedian

import "container/heap"

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
	e := x.(*item)
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

type MovingMedian struct {
	queueIndex int
	nitems     int
	queue      []item
	maxHeap    maxItemHeap
	minHeap    minItemHeap
}

func NewMovingMedian(size int) MovingMedian {
	m := MovingMedian{
		queue:   make([]item, size),
		maxHeap: maxItemHeap{},
		minHeap: minItemHeap{},
	}

	heap.Init(&m.maxHeap)
	heap.Init(&m.minHeap)
	return m
}

func (m *MovingMedian) Push(v float64) {
	if len(m.queue) == 1 {
		m.queue[0].f = v
		return
	}

	itemPtr := &m.queue[m.queueIndex]
	m.queueIndex++
	if m.queueIndex >= len(m.queue) {
		m.queueIndex = 0
	}

	minHeapLen := m.minHeap.Len()
	if m.nitems == len(m.queue) {
		if itemPtr.heapIndex < minHeapLen && itemPtr == m.minHeap.itemHeap[itemPtr.heapIndex] {
			if v >= m.maxHeap.itemHeap[0].f {
				itemPtr.f = v
				heap.Fix(&m.minHeap, itemPtr.heapIndex)
				return
			}

			moveItem := m.maxHeap.itemHeap[0]
			moveItem.heapIndex = itemPtr.heapIndex
			m.minHeap.itemHeap[itemPtr.heapIndex] = moveItem
			itemPtr.f = v
			m.maxHeap.itemHeap[0] = itemPtr

			heap.Fix(&m.minHeap, itemPtr.heapIndex)
			itemPtr.heapIndex = 0
			heap.Fix(&m.maxHeap, 0)
			return
		} else {
			if v <= m.minHeap.itemHeap[0].f {
				itemPtr.f = v
				heap.Fix(&m.maxHeap, itemPtr.heapIndex)
				return
			}

			moveItem := m.minHeap.itemHeap[0]
			moveItem.heapIndex = itemPtr.heapIndex
			m.maxHeap.itemHeap[itemPtr.heapIndex] = moveItem
			itemPtr.f = v
			m.minHeap.itemHeap[0] = itemPtr

			heap.Fix(&m.maxHeap, itemPtr.heapIndex)
			itemPtr.heapIndex = 0
			heap.Fix(&m.minHeap, 0)
			return
		}
	}

	m.nitems++
	itemPtr.f = v
	if minHeapLen == 0 {
		heap.Push(&m.minHeap, itemPtr)
	} else if v > m.minHeap.itemHeap[0].f {
		heap.Push(&m.minHeap, itemPtr)
		if minHeapLen > m.maxHeap.Len() {
			moveItem := heap.Pop(&m.minHeap)
			heap.Push(&m.maxHeap, moveItem)
		}
	} else {
		heap.Push(&m.maxHeap, itemPtr)
		if m.maxHeap.Len() == (minHeapLen + 2) {
			moveItem := heap.Pop(&m.maxHeap)
			heap.Push(&m.minHeap, moveItem)
		}
	}
}

func (m *MovingMedian) Median() float64 {
	if len(m.queue) == 1 {
		return m.queue[0].f
	}

	if m.maxHeap.Len() == m.minHeap.Len() {
		return (m.maxHeap.itemHeap[0].f + m.minHeap.itemHeap[0].f) / 2
	}

	if m.maxHeap.Len() > m.minHeap.Len() {
		return m.maxHeap.itemHeap[0].f
	}

	return m.minHeap.itemHeap[0].f
}
