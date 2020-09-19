package hw1

import (
	"math"
	"testing"
)

func TestReverse(t *testing.T) {
	cases := []float64{2, 3, 5, 7, 9, 16, 256, 1024, 4096}
	for _, c := range cases {
		result := Sqrt(c)
		if result-math.Sqrt(c) > 0.01 || math.Sqrt(c)-result > 0.01 {
			t.Errorf("my Sqrt(%f) == %f, but want %f", c, result, math.Sqrt(c))
		}
	}
}
