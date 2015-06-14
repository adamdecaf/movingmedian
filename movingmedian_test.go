package movingmedian

import (
	"github.com/wangjohn/quickselect"
	"log"
	"math"
	"math/rand"
	"testing"
)

func TestUnit(t *testing.T) {
	tests := []struct {
		name       string
		windowSize int
		data       []float64
		want       []float64
	}{
		{
			"OneWindowSize",
			1,
			[]float64{1, 3, 5, 7, 9, 11, math.NaN()},
			[]float64{1, 3, 5, 7, 9, 11, math.NaN()},
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
			"DecreasingIncreasingValues",
			4,
			[]float64{21, 19, 17, 15, 13, 11, 13, 15, 17, 19},
			[]float64{21, 20, 19, 18, 16, 14, 13, 13, 14, 16},
		},
		{
			"IncreasingDecreasingValues",
			4,
			[]float64{11, 13, 15, 17, 19, 21, 19, 17, 15, 13},
			[]float64{11, 12, 13, 14, 16, 18, 19, 19, 18, 16},
		},
		{

			"ZigZag",
			4,
			[]float64{21, 23, 17, 27, 13, 31, 9, 35, 5, 39, 1},
			[]float64{21, 22, 21, 22, 20, 22, 20, 22, 20, 22, 20},
		},
		{

			"NewValuesInBetween",
			4,
			[]float64{21, 21, 19, 19, 21, 21, 19, 19, 19, 19},
			[]float64{21, 21, 21, 20, 20, 20, 20, 20, 19, 19},
		},
		{
			"SameNumberInBothHeaps3Times",
			4,
			[]float64{11, 13, 13, 13, 25, 27, 29, 31},
			[]float64{11, 12, 13, 13, 13, 19, 26, 28},
		},
		{
			"SameNumberInBothHeaps3TimesDecreasing",
			4,
			[]float64{31, 29, 29, 29, 17, 15, 13, 11},
			[]float64{31, 30, 29, 29, 29, 23, 16, 14},
		},
		{
			"SameNumberInBothHeaps4Times",
			4,
			[]float64{11, 13, 13, 13, 13, 25, 27, 29, 31},
			[]float64{11, 12, 13, 13, 13, 13, 19, 26, 28},
		},
	}

	for _, test := range tests {
		log.Println("test name", test.name)
		m := NewMovingMedian(test.windowSize)
		for i, v := range test.data {
			m.Push(v)
			actual := m.Median()
			if test.want[i] != actual && !(math.IsNaN(actual) && math.IsNaN(test.want[i])) {
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

func TestRandom(t *testing.T) {
	rangeSize := 100
	for windowSize := 1; windowSize < 5; windowSize++ {
		data := getData(rangeSize, windowSize)
		intData := make([]int, rangeSize)
		for i, v := range data {
			intData[i] = int(v)
		}

		log.Println("test name random test window size", windowSize)
		m := NewMovingMedian(windowSize)
		for i, v := range data {
			want := median(data, i, windowSize)

			m.Push(v)
			actual := m.Median()
			if want != actual {
				firstElement := 1 + i - windowSize
				if firstElement < 0 {
					firstElement = 0
				}

				t.Errorf("failed on test random window size %d index %d the median of %d is %f and not %f",
					windowSize,
					i,
					intData[firstElement:1+i],
					want,
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
	data := getData(numberOfValues, windowSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := NewMovingMedian(windowSize)
		for _, v := range data {
			m.Push(v)
			m.Median()
		}
	}
}

func getData(rangeSize, windowSize int) []float64 {
	var data = make([]float64, rangeSize)
	var r = rand.New(rand.NewSource(99))
	for i, _ := range data {
		data[i] = r.Float64()
	}

	return data
}

func median(data []float64, i, windowSize int) float64 {
	min := 1 + i - windowSize
	if min < 0 {
		min = 0
	}

	window := make([]float64, 1+i-min)
	copy(window, data[min:i+1])
	return percentile(window, 50, true)
}

func percentile(data []float64, percent float64, interpolate bool) float64 {
	if len(data) == 0 || percent < 0 || percent > 100 {
		return math.NaN()
	}
	if len(data) == 1 {
		return data[0]
	}

	k := (float64(len(data)-1) * percent) / 100
	length := int(math.Ceil(k)) + 1
	quickselect.Float64QuickSelect(data, length)
	top, secondTop := math.Inf(-1), math.Inf(-1)
	for _, val := range data[0:length] {
		if val > top {
			secondTop = top
			top = val
		} else if val > secondTop {
			secondTop = val
		}
	}
	remainder := k - float64(int(k))
	if remainder == 0 || !interpolate {
		return top
	}
	return (top * remainder) + (secondTop * (1 - remainder))
}
