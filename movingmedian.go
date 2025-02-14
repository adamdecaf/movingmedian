// Package movingmedian computes the median of a windowed stream of data.
package movingmedian

import "sort"

// MovingMedian computes the median over a fixed-size sliding window.
// It uses a ring buffer to track insertion order and a sorted slice for median queries.
type MovingMedian struct {
	size   int       // fixed window size
	ring   []float64 // ring buffer in insertion order
	sorted []float64 // sorted copy of current window
	pos    int       // next position in ring buffer to write
	count  int       // number of elements currently in the window (<= size)
}

// NewMovingMedian creates a new MovingMedian for a window of the given size.
func NewMovingMedian(size int) *MovingMedian {
	return &MovingMedian{
		size:   size,
		ring:   make([]float64, size),
		sorted: make([]float64, 0, size),
	}
}

// Push adds a new observation. If the window is full, it removes the oldest.
func (m *MovingMedian) Push(v float64) {
	if m.count < m.size {
		// Window not yet full: add value to ring buffer and sorted slice.
		m.ring[m.pos] = v
		m.pos = (m.pos + 1) % m.size

		// Insert v into the sorted slice.
		i := sort.Search(m.count, func(i int) bool {
			return m.sorted[i] >= v
		})
		m.sorted = append(m.sorted, 0)     // extend slice by one
		copy(m.sorted[i+1:], m.sorted[i:]) // shift right
		m.sorted[i] = v
		m.count++
	} else {
		// Window is full: remove the oldest value then insert the new one.
		old := m.ring[m.pos]
		m.ring[m.pos] = v
		m.pos = (m.pos + 1) % m.size

		// Remove one occurrence of old from the sorted slice.
		i := sort.Search(m.count, func(i int) bool {
			return m.sorted[i] >= old
		})
		// There may be duplicates; we expect sorted[i]==old.
		if i < m.count && m.sorted[i] == old {
			copy(m.sorted[i:], m.sorted[i+1:])
			m.sorted = m.sorted[:m.count-1]
		} else {
			// Should never happen.
			panic("old element not found in sorted slice")
		}

		// Insert new value v into the sorted slice.
		j := sort.Search(m.count-1, func(i int) bool {
			return m.sorted[i] >= v
		})
		m.sorted = append(m.sorted, 0)
		copy(m.sorted[j+1:], m.sorted[j:])
		m.sorted[j] = v
	}
}

// Median returns the current median of the window.
// For an odd count, it returns the middle element; for an even count, the average of the two.
func (m *MovingMedian) Median() float64 {
	if m.count == 0 {
		return 0.0
	}
	if m.count%2 == 1 {
		return m.sorted[m.count/2]
	}
	// For even count, average the two middle values.
	return (m.sorted[m.count/2-1] + m.sorted[m.count/2]) / 2.0
}
