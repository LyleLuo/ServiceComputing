package marshal

import (
	"encoding/json"
	"testing"
)

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
	Rate   float32 `json:"omitempty"`
	actor  string
	People Watcher
}

type Activity struct {
	ID     int
	Movies Movie
	Nation
}

// 综合测试
func TestJsonMarshal(t *testing.T) {
	a := Activity{18342069, Movie{2020, false, "Xiyangyang", 4.5, "Pan Maolin", Watcher{"Luo Weile", "Male"}}, Nation{"China"}}
	data, err := JsonMarshal(a)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(a)
	if string(data) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(data))
	}
}
