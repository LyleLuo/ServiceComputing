package hw2

import "testing"

func TestSort(t *testing.T) {
	arr := []int{3, 4, 2, 1, 7, 5, 6, 8, 9, 0}
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	Sort(arr, 0, len(arr)-1)
	for i := 0; i < 10; i++ {
		if arr[i] != expected[i] {
			t.Errorf("expected '%v' but got '%v'", expected, arr)
			break
		}
	}
}

func BenchmarkSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := []int{3, 4, 2, 1, 2345, 4325, 6, 785, 875, 958, 46, 2435, 6, 76, 98, 897, 576, 28, 82, 7, 5, 6, 8, 9, 0}
		Sort(arr, 0, len(arr)-1)
	}
}
