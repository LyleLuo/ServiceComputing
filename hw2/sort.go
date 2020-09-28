package hw2

func swap(a, b *int) {
	*a, *b = *b, *a
}

func partition(arr []int, low, high int) int {
	lastSmall := low
	swap(&arr[low], &arr[(low+high)/2])
	pivot := arr[low]
	for i := low + 1; i < high+1; i++ {
		if arr[i] < pivot {
			lastSmall++
			swap(&arr[lastSmall], &arr[i])
		}
	}
	swap(&arr[lastSmall], &arr[low])
	return lastSmall
}

// Sort 自己写的快速排序
func Sort(arr []int, low, high int) {
	if low < high {
		pivotPos := partition(arr, low, high)
		Sort(arr, low, pivotPos-1)
		Sort(arr, pivotPos+1, high)
	}
}
