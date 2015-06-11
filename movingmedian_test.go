package movingmedian

import (
	"math/rand"
	"testing"
)

func TestSameNumberInBothHeaps(t *testing.T) {
	tests := []struct {
		name       string
		windowSize int
		data       []float64
		want       []float64
	}{
		{
			"OneWindowSize",
			1,
			[]float64{1, 3, 5, 7, 9, 11},
			[]float64{1, 3, 5, 7, 9, 11},
		},

		{
			"OddWindowSize",
			3,
			[]float64{1, 3, 5, 7, 9, 11},
			[]float64{1, 2, 3, 5, 7, 9},
		},
		{
			"EvenWindowSize",
			4,
			[]float64{1, 3, 5, 7, 9, 11},
			[]float64{1, 2, 3, 4, 6, 8},
		},
		{
			"DecreasingValues",
			4,
			[]float64{19, 17, 15, 13, 11, 9},
			[]float64{19, 18, 17, 16, 14, 12},
		},
		{
			"SameNumberInBothHeaps3Times",
			4,
			[]float64{11, 13, 13, 13, 25, 27, 29, 31},
			[]float64{11, 12, 13, 13, 13, 19, 26, 28},
		},
		{
			"SameNumberInBothHeaps4Times",
			4,
			[]float64{11, 13, 13, 13, 13, 25, 27, 29, 31},
			[]float64{11, 12, 13, 13, 13, 13, 19, 26, 28},
		},
	}

	for _, test := range tests {
		m := NewMovingMedian(test.windowSize)
		for i, v := range test.data {
			m.Push(v)
			actual := m.Median()
			if test.want[i] != actual {
				firstElement := 1 + i - test.windowSize
				if firstElement < 0 {
					firstElement = 0
				}
				t.Errorf("failed on test %s index %d the median of %f is %f and not %f",
					test.name,
					i,
					test.data[firstElement:1+i],
					test.want[i],
					actual)
			}
		}
	}
}

func Benchmark_10values_windowsize1(b *testing.B) {
	benchmark(b, 10, 1)
}

func Benchmark_100values_windowsize10(b *testing.B) {
	benchmark(b, 100, 10)
}

func Benchmark_10Kvalues_windowsize100(b *testing.B) {
	benchmark(b, 10000, 100)
}

func Benchmark_10Kvalues_windowsize1000(b *testing.B) {
	benchmark(b, 10000, 1000)
}

func benchmark(b *testing.B, numberOfValues, windowSize int) {
	data := getData(numberOfValues)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewMovingMedian(windowSize)
		for _, v := range data {
			m.Push(v)
			m.Median()
		}
	}
}

func getData(rangeSize int) []float64 {
	var data = make([]float64, rangeSize)
	var r = rand.New(rand.NewSource(99))
	for i, _ := range data {
		data[i] = r.Float64()
	}

	return data
}
