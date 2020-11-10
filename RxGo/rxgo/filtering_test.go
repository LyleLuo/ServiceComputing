package rxgo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/LyleLuo/ServiceComputing/RxGo/rxgo"
)

func TestDebounce(t *testing.T) {
	answer := []int{3, 9}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 && x != 4 {
			time.Sleep(10 * time.Millisecond)
		} else if x == 4 {
			time.Sleep(30 * time.Millisecond)
		}
		return x
	}).Debounce(20 * time.Millisecond).Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestDistinct(t *testing.T) {
	answer := []int{0, 1, 2, 3, 4, 5}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 4, 1).Distinct().Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestElementAt(t *testing.T) {
	var result int
	answer := 1
	rxgo.Just(0, 1, 2, 3, 9).ElementAt(1).Subscribe(func(x int) {
		result = x
	})
	if result != answer {
		t.Errorf("Excepted number: %v, but your number: %v\n", answer, result)
	}
}

func TestFirst(t *testing.T) {
	var result int
	answer := -1
	rxgo.Just(-1, 1, 2, 3).First().Subscribe(func(x int) {
		result = x
	})
	if result != answer {
		t.Errorf("Excepted number: %v, but your number: %v\n", answer, result)
	}
}

func TestIgnoreElements(t *testing.T) {
	var result []int
	rxgo.Just(0, 1, 2, 3, 9).IgnoreElements().Subscribe(func(x int) {
		result = append(result, x)
	})
	if result != nil {
		t.Errorf("Ignore elements function occured error!\n")
	}
}

func TestLast(t *testing.T) {
	var result int
	answer := 9
	rxgo.Just(6, 2, 3, 9).Last().Subscribe(func(x int) {
		result = x
	})
	if result != answer {
		t.Errorf("Excepted number: %v, but your number: %v", answer, result)
	}
}

func TestSample(t *testing.T) {
	answer := []int{2, 4, 6, 8, 9}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 {
			time.Sleep(100 * time.Millisecond)
		}
		return x
	}).Sample(210 * time.Millisecond).Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestSkip(t *testing.T) {
	answer := []int{4, 5, 6, 7, 8, 9}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Skip(4).Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestSkiplast(t *testing.T) {
	answer := []int{0, 1, 2, 3}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Skiplast(6).Subscribe(func(x int) {
		result = append(result, x)
	})
	println(answer)
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestTake(t *testing.T) {
	answer := []int{0, 1, 2, 3}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Take(4).Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}

func TestTakelast(t *testing.T) {
	answer := []int{7, 8, 9}
	var result []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Takelast(3).Subscribe(func(x int) {
		result = append(result, x)
	})
	if !reflect.DeepEqual(answer, result) {
		t.Errorf("Excepted slice: %v, but your slice: %v\n", answer, result)
	}
}
