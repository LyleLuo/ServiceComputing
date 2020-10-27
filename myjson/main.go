package main

import (
	"fmt"
	"log"

	"github.com/LyleLuo/ServiceComputing/myjson/marshal"
)

func main() {
	type Watcher struct {
		PeopleName string
		sex        string
	}

	type Nation struct {
		NationName string
	}

	type Movie struct {
		Year   int  `json:"released"`
		Color  bool `json:"color"`
		Name   string
		Rate   float32 `json:"rate,omitempty"`
		actor  string
		People Watcher
	}

	type Activity struct {
		ID     int
		Movies Movie
		Nation
	}
	// 可以传入各种（嵌套）结构体（数组）或者map
	a := [3]Activity{
		Activity{18342069, Movie{2020, false, "Xiyangyang", 4.5, "Pan Maolin", Watcher{"Luo Weile", "Male"}}, Nation{"China"}},
		Activity{12345678, Movie{2020, false, "Huitailang", 5.0, "Luo Jun", Watcher{"Luo Luo", "Female"}}, Nation{"USA"}},
		Activity{98765433, Movie{2020, false, "Xiongchumo", 4.3, "Dilireba", Watcher{"pml", "Male"}}, Nation{"UK"}}}

	data, err := marshal.JsonMarshal(a)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}
