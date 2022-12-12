package test

import (
	"math"
	"testing"
)

func Cosine(a []float64, b []float64) float64 {
	var (
		aLen  = len(a)
		bLen  = len(b)
		s     = 0.0
		sa    = 0.0
		sb    = 0.0
		count = 0
	)
	if aLen > bLen {
		count = aLen
	} else {
		count = bLen
	}
	for i := 0; i < count; i++ {
		if i >= bLen {
			sa += math.Pow(a[i], 2)
			continue
		}
		if i >= aLen {
			sb += math.Pow(b[i], 2)
			continue
		}
		s += a[i] * b[i]
		sa += math.Pow(a[i], 2)
		sb += math.Pow(b[i], 2)
	}
	return s / (math.Sqrt(sa) * math.Sqrt(sb))
}
func Test_cosin_simulation(t *testing.T) {
	t.Log(Cosine([]float64{1, 2, 3, 4, 5, 6}, []float64{3, 2, 1}))
	t.Log(Cosine([]float64{1, 2, 3}, []float64{4, 4, 5, 6}))
	t.Log(Cosine([]float64{1, 1, 1}, []float64{1, 2, 3}))
	t.Log(Cosine([]float64{1, 1, 1}, []float64{1, 1, 1}))
}
