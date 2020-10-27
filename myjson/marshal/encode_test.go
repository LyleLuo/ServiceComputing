package marshal

import (
	"encoding/json"
	"testing"
)

type testString struct {
	S1 string
	S2 string
}

type testFloat struct {
	F32 float32
	F64 float64
}

type testInt struct {
	I32 int32
	U32 uint32
	I64 int64
	U64 uint64
}

func TestEncodeStrOfStruct(t *testing.T) {
	s := testString{"Hello", "world"}
	e := &EncodeState{}
	err := e.marshal(s)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(s)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}

func TestEncodeIntOfStruct(t *testing.T) {
	s := testInt{1342543, 243524432, -5445646134, 24562463546456536}
	e := &EncodeState{}
	err := e.marshal(s)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(s)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}

func TestEncodeFloatOfStruct(t *testing.T) {
	s := testFloat{1.5, 3.14}
	e := &EncodeState{}
	err := e.marshal(s)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(s)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}
func TestEncodeMapIntStr(t *testing.T) {
	m := map[int]string{1: "hello", 2: "world"}
	e := &EncodeState{}
	err := e.marshal(m)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(m)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}
func TestEncodeMapIntFloat(t *testing.T) {
	m := map[int]float32{1: 43252.2, 2: 3.14}
	e := &EncodeState{}
	err := e.marshal(m)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(m)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}

func TestEncodeMapStrInt(t *testing.T) {
	m := map[string]int{"lwl": 666, "pml": 520}
	e := &EncodeState{}
	err := e.marshal(m)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(m)

	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}

func TestEncodeStructArray(t *testing.T) {
	s := [3]testString{
		{"Hello", "world"},
		{"Hi", "pml"},
		{"pml", "666"}}
	e := &EncodeState{}
	err := e.marshal(s)
	if err != nil {
		t.Error(err)
	}
	expect, _ := json.Marshal(s)
	if string(e.Bytes()) != string(expect) {
		t.Errorf("expect string is %s, but output is %s", expect, string(e.Bytes()))
	}
}
