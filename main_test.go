package main

import (
	"github.com/wangjohn/quickselect"
	"math"
	"math/rand"
	"testing"
)

func TestMedian(t *testing.T) {
	var windowSize = 10
	data := getData(100)
	m := NewMovingMedian(windowSize)
	for i, v := range data {
		want := median(data, i, windowSize)

		m.Push(v)
		actual := m.Median()
		if want != actual {
			t.Errorf("median failed on index %v: item %v want %v actual %v", i, v, want, actual)
		}
	}
}

func BenchmarkMovingMedianOptimized(b *testing.B) {
	var windowSize = int(4e2)
	data := getData(int(16e4))

	m := NewMovingMedian(windowSize)
	for _, v := range data {
		m.Push(v)
		m.Median()
	}
}

func BenchmarkMovingMedian(b *testing.B) {
	var windowSize = int(4e2)
	data := getData(int(16e4))

	for i := range data {
		want := median(data, i, windowSize)
		want = want
	}
}

func getData(rangeSize int) []float64 {
	var data = make([]float64, rangeSize)
	var r = rand.New(rand.NewSource(99))
	for i, _ := range data {
		data[i] = float64(int(10 * r.Float64()))
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
