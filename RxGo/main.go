package main

import (
	"fmt"
	"time"

	"github.com/LyleLuo/ServiceComputing/RxGo/rxgo"
)

func main() {
	var result1 []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 && x != 4 {
			time.Sleep(10 * time.Millisecond)
		} else if x == 4 {
			time.Sleep(30 * time.Millisecond)
		}
		return x
	}).Debounce(20 * time.Millisecond).Subscribe(func(x int) {
		result1 = append(result1, x)
	})
	fmt.Println(result1)

	var result2 []string
	rxgo.Just("rxgo", "rxgo", "rxgo").Distinct().Subscribe(func(x string) {
		result2 = append(result2, x)
	})
	fmt.Println(result2)

	var result3 string
	rxgo.Just("rxgo", "rxgo", "rxgo").ElementAt(1).Subscribe(func(x string) {
		result3 = x
	})
	fmt.Println(result3)

	var result4 int
	rxgo.Just(-1, 1, 2, 3).First().Subscribe(func(x int) {
		result4 = x
	})
	fmt.Println(result4)

	var result5 []int
	rxgo.Just(0, 1, 2, 3, 9).IgnoreElements().Subscribe(func(x int) {
		result5 = append(result5, x)
	})
	fmt.Println(result5)

	var result6 string
	rxgo.Just("pml", "lwl").Last().Subscribe(func(x string) {
		result6 = x
	})
	fmt.Println(result6)

	var result7 []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Map(func(x int) int {
		if x != 0 {
			time.Sleep(100 * time.Millisecond)
		}
		return x
	}).Sample(210 * time.Millisecond).Subscribe(func(x int) {
		result7 = append(result7, x)
	})
	fmt.Println(result7)

	var result8 []string
	rxgo.Just("0", "1", "2", "3", "4", "5", "6", "7", "8", "9").Skip(4).Subscribe(func(x string) {
		result8 = append(result8, x)
	})
	fmt.Println(result8)

	var result9 []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Skiplast(6).Subscribe(func(x int) {
		result9 = append(result9, x)
	})
	fmt.Println(result9)

	var result10 []string
	rxgo.Just("i", "love", "you", "TA").Take(4).Subscribe(func(x string) {
		result10 = append(result10, x)
	})
	fmt.Println(result10)

	var result11 []int
	rxgo.Just(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Takelast(3).Subscribe(func(x int) {
		result11 = append(result11, x)
	})
	fmt.Println(result11)
}
