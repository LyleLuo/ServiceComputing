# 作业二：使用快速排序练习TDD
## 先写测试
```go
package hw2

import "testing"

func TestSort(t *testing.T) {
	arr := []int{3, 4, 2, 1, 7, 5, 6, 8, 9, 0}
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

    Sort(arr, 0, len(arr) - 1)
	for i := 0; i < 10; i++ {
		if arr[i] != expected[i] {
            t.Errorf("expected '%v' but got '%v'", expected, arr)
            break
		}
	}
}
```

## 尝试运行测试
```
[luowle@VM_0_4_centos hw2]$ go test
# github.com/ServiceComputing/hw2 [github.com/ServiceComputing/hw2.test]
./sort_test.go:9:2: undefined: Sort
FAIL    github.com/ServiceComputing/hw2 [build failed]
```

## 先使用最少的代码来让失败的测试先跑起来
```go
package hw2

// Sort 自己写的快速排序
func Sort(arr []int, start, end int) {
}

```
go test的结果如下，因为Sort函数没有任何内容，因此不能通过测试
```
[luowle@VM_0_4_centos hw2]$ go test
--- FAIL: TestSort (0.00s)
    sort_test.go:12: expected '[0 1 2 3 4 5 6 7 8 9]' but got '[3 4 2 1 7 5 6 8 9 0]'
FAIL
exit status 1
FAIL    github.com/ServiceComputing/hw2 0.003s
```
## 把代码补充完整，使得它能够通过测试
```go
package hw2

func partition(arr []int, low, high int) int {
	lastSmall := low
	arr[low], arr[(low+high)/2] = arr[(low+high)/2], arr[low]
	pivot := arr[low]
	for i := low + 1; i < high+1; i++ {
		if arr[i] < pivot {
			lastSmall++
			arr[lastSmall], arr[i] = arr[i], arr[lastSmall]
		}
	}
	arr[lastSmall], arr[low] = arr[low], arr[lastSmall]
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
```
测试结果如下
```
[luowle@VM_0_4_centos hw2]$ go test
PASS
ok      github.com/ServiceComputing/hw2 0.002s
```

## 重构
通过写一个swap函数让代码看起来更有逻辑，完整代码如下所示
```go
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

```

## 基准测试
编写一个基准测试如下所示
```go
func BenchmarkSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := []int{3, 4, 2, 1, 2345, 4325, 6, 785, 875, 958, 46, 2435, 6, 76, 98, 897, 576, 28, 82, 7, 5, 6, 8, 9, 0}
		Sort(arr, 0, len(arr)-1)
	}
}
```

运行测试结果如下所示，每次排序大概有300-400ns
```
[luowle@VM_0_4_centos hw2]$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/ServiceComputing/hw2
BenchmarkSort    3297258               342 ns/op
PASS
ok      github.com/ServiceComputing/hw2 1.509s
```

## 总结
- 理解TDD、重构、测试、基准测试等概念。
- 熟悉了Go语言的语法
- 加深了工作空间组织、包、变量、控制、函数、结构、集合等的认识