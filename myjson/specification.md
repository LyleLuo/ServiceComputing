# specification
<!-- TOC -->

- [specification](#specification)
    - [设计说明](#设计说明)
        - [设计思路](#设计思路)
        - [文件](#文件)
        - [核心功能 JsonMarshal](#核心功能-jsonmarshal)
        - [函数返回值](#函数返回值)
        - [核心的设计与实现](#核心的设计与实现)
            - [func typeEncoder(t reflect.Type) EncoderFunc](#func-typeencodert-reflecttype-encoderfunc)
            - [func newTypeEncoder(t reflect.Type, allowAddr bool) EncoderFunc](#func-newtypeencodert-reflecttype-allowaddr-bool-encoderfunc)
            - [func intEncoder(e *EncodeState, v reflect.Value)](#func-intencodere-encodestate-v-reflectvalue)
            - [func (bits floatEncoder) encode(e *EncodeState, v reflect.Value)](#func-bits-floatencoder-encodee-encodestate-v-reflectvalue)
            - [func newStructEncoder(t reflect.Type) EncoderFunc](#func-newstructencodert-reflecttype-encoderfunc)
            - [func typeFields(t reflect.Type) structFields](#func-typefieldst-reflecttype-structfields)
            - [func (se structEncoder) encode(e *EncodeState, v reflect.Value)](#func-se-structencoder-encodee-encodestate-v-reflectvalue)
            - [func interfaceEncoder(e *EncodeState, v reflect.Value)](#func-interfaceencodere-encodestate-v-reflectvalue)
    - [单元或集成测试](#单元或集成测试)
        - [测试代码](#测试代码)
            - [marshal_test.go](#marshal_testgo)
            - [encode_test.go](#encode_testgo)
            - [tags_test.go](#tags_testgo)
        - [测试结果](#测试结果)
    - [功能测试](#功能测试)
        - [测试代码](#测试代码-1)
        - [测试结果](#测试结果-1)
    - [总结](#总结)

<!-- /TOC -->
## 设计说明
### 设计思路
本次作业的实现是~~当代码人肉翻译器~~参考官方库的实现，实现了一个**精简版**的结构数据->JSON解析器。
### 文件
- marshal.go -- 实现核心函数 JsonMarshal
- encode.go -- 实现了解析 JSON 的具体实现
- tags.go -- 解析结构体内 tag 所用到的函数
- xx_test.go -- xx对应的测试文件

### 核心功能 JsonMarshal
JsonMarshal 用于将结构数据解析成为字节流，支持的类型包括（嵌套）结构体，map，结构体数组。结构内的数据必须是 go 的基本类型
```go
// JsonMarshal 用于将结构数据解析成为字节流，支持的类型包括（嵌套）结构体，map，结构体数组。结构内的数据必须是go的基本类型
func JsonMarshal(v interface{}) ([]byte, error) {
	//存储v的编组之后的bytes.Buffer对象（见下 具体数据结构）
	e := &EncodeState{}
	//调用编组对象的marshal方法，encOpts 对象：编码过程的配置
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}
```

### 函数返回值
函数会返回两个值
1. 第一个返回值 []byte 返回的是解析后的JSON内容，如果解析过程中出错则会返回 nil
2. 第二个返回值 error 返回的是解析过程中出错的内容，如果没有错误则会返回nil

### 核心的设计与实现
1. 首先，经过函数调用到达 typeEncoder，查看缓存是否存在该类型的解析函数。
2. 若没有，则使用 newTypeEncoder 创建新的解析函数用于缓存。
3. newTypeEncoder 会通过 reflect.kind 获得类型，从而返回已实现类型的解析函数。
4. 如果是 struct、map 等类型则要为每个 reflect.type 使用相应方法创建解析函数。
5. 为 struct 创建解析方法要先使用 typeFields 解析结构体内的（匿名）变量，层次等。
6. 然后使用 reflect.type 专属的解析函数进行解析。
7. 结构体/ map /interface 在解析过程中会调用其成员/键值对相应的解析函数进行解析并写到缓冲区中。


#### func typeEncoder(t reflect.Type) EncoderFunc
该函数作为整个程序的核心，根据输入的函数类型来返回一个合适的解析函数。其中，参考官方文档使用了一个map作为缓存，使得已经使用过的 type 不需要再重新生成 EncoderFunc，从而提高性能。操作步骤如下：
1. 尝试在缓存中读取 EncoderFunc，如果成功则直接返回
2. Cache对象里没有，则自己创建
3. 将创建好的 EncoderFunc 对象以type:encoderFunc 的形式存入缓存 map 中
4. 返回新 EncoderFunc

代码如下
```go
// EncoderFunc 将reflect.Value写到EncodeState的buffer函数
type EncoderFunc func(e *EncodeState, v reflect.Value)

func typeEncoder(t reflect.Type) EncoderFunc {
	//encoderCache 是个sync.Map 缓存处理的encoderFunc函数，如果能找到，直接使用
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(EncoderFunc)
	}
	var (
		wg sync.WaitGroup
		f  EncoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, EncoderFunc(func(e *EncodeState, v reflect.Value) {
		wg.Wait()
		f(e, v)
	}))
	if loaded {
		return fi.(EncoderFunc)
	}
	//Cache对象里没有，则自己创建
	f = newTypeEncoder(t, true)
	wg.Done()
	//将创建好的 EncoderFunc 对象以type：encoderFunc 的形式存入缓存map中
	encoderCache.Store(t, f)
	return f
}
```

#### func newTypeEncoder(t reflect.Type, allowAddr bool) EncoderFunc
在上一个功能中，用该函数新建一个对应类型的 EncoderFunc。对于 go 自带的基本类型，可以返回同一个处理函数。但是对于 struct、map 和 interface，就要针对不同的 reflect.Type 生成不同的 EncoderFunc 供于缓存。看源码可清晰知道我的实现所支持的类型，代码如下：
```go
// 当需要用的处理函数没有被缓存时，则新建一个并加进缓存
func newTypeEncoder(t reflect.Type, allowAddr bool) EncoderFunc {
	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}
```
下面描述几个典型的 EncoderFunc

#### func intEncoder(e *EncodeState, v reflect.Value)
这个函数用于处理 reflect.Int 大类的。EncodeState 中有缓冲区，Write 是 EncodeState 的一个方法，用于将 []byte 写到 EncodeState 的缓冲区中。直接调用 strconv.AppendInt 就可以将 int 写到字节流中。
```go
// 用于处理int和uint的encoderFunc
func intEncoder(e *EncodeState, v reflect.Value) {
	b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
	e.Write(b)
}
```

#### func (bits floatEncoder) encode(e *EncodeState, v reflect.Value)
上面 newTypeEncoder 所选择的 float32Encoder 被定义为 (floatEncoder(32)).encode。下面是该方法的实现。该方法通过数据范围，来决定是否采用科学计数法写到缓冲区上。
```go
// 用于解析float32和float64的方法
func (bits floatEncoder) encode(e *EncodeState, v reflect.Value) {
	f := v.Float()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		e.error(&unsupportedValueError{v, strconv.FormatFloat(f, 'g', -1, int(bits))})
	}

	// 以下决定是否使用科学计数法
	b := e.scratch[:0]
	abs := math.Abs(f)
	fmt := byte('f')
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
	if fmt == 'e' {
		// 将 e-09 设为 e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}

	e.Write(b)
}
```

#### func newStructEncoder(t reflect.Type) EncoderFunc
新建一个解析器，对于结构体来说，每一种结构体都要新建一个解析方法用于缓存。

cachedTypeFields 和上面使用缓存的方法类似。如果缓存不存在，则使用 **typeFields** 函数获得一个新的 structFields。structFields 记录着结构体的各种信息，之后会返回一个 **structFields.encode** 的方法用于解析结构体，该方法会在后面讲述设计。
```go
// 缓存 field 函数从而避免重复操作
func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}
// 新建一个解析器，对于结构体来说，每一种结构体都要新建一个解析方法用于缓存
func newStructEncoder(t reflect.Type) EncoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}
```
#### func typeFields(t reflect.Type) structFields
解析结构体所需的类型如下所示
```go
type structEncoder struct {
	fields structFields
}

type structFields struct {
	list      []field
	nameIndex map[string]int
}

// A field represents a single field found in a struct.
type field struct {
	name      string
	nameBytes []byte // []byte(name)

	nameNonEsc string // `"` + name + `":`

	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool

	encoder EncoderFunc
}
```
该函数的目的就是将结构体的内容解析到 fields 中，具体的代码可以直接查看文件，文件内有一定的解释用以辅助阅读代码。思路如下：

1. 从当前结构体开始，遍历成员变量
2. 当前层的 len(index) 是 1，如果有匿名结构体，就把该匿名结构体添加进队列中。待该层遍历完则出队分析。
3. 第 N 层的 len(index) 是 N，该匿名层的成员变量将会被第一层使用，每个匿名结构体成员变量的位置会 append 到 index 中。（index 十分重要，encode 匿名结构体的时候发挥重要作用）
4. 遍历的过程中分析 name、tag 等信息

由于篇幅所限，在这里仅展示部分重要设计。

- 跳过匿名且未导出的非结构体类型
```go
if sf.Anonymous {
	t := sf.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if isUnexported && t.Kind() != reflect.Struct {
		continue
	}
} else if isUnexported {
	continue
}
```
- 解析 tag
```go
tag := sf.Tag.Get("json")
if tag == "-" {
	continue
}
name, opts := parseTag(tag)
if !isValidTag(name) {
	name = ""
}
```
#### func (se structEncoder) encode(e *EncodeState, v reflect.Value)
该方法是用于解析结构体的方法。
- 此方法依赖于上面一个函数得到的fields。
- 该算法的核心是利用上面的 index 。当 list[i] 所存的是匿名结构体中的变量时，需要借助 index 进入到 len(index) 层中读取 value。
- 而把键值对输出到缓冲区的的操作看代码就能懂，此处不再赘述。代码如下：

```go
func (se structEncoder) encode(e *EncodeState, v reflect.Value) {
	next := byte('{')
FieldLoop:
	for i := range se.fields.list {
		f := &se.fields.list[i]

		// Find the nested struct field by following f.index.
		// 用数组实现的递归搜索
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue FieldLoop
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}
		
		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		e.WriteByte(next)
		next = ','
		e.WriteString(f.nameNonEsc)

		f.encoder(e, fv)
	}
	if next == '{' {
		e.WriteString("{}")
	} else {
		e.WriteByte('}')
	}
}
```

#### func interfaceEncoder(e *EncodeState, v reflect.Value)
无论是array、map、还是 interface，本质上都是拆开里面的类型然后使用上面的 typeEncoder 获得解析函数从而进行使用。下面是以 interface 为例子。
```go
// 用于解析interface类型的函数，本质上是获得接口的类型后在调用相应函数
func interfaceEncoder(e *EncodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// Elem返回接口v包含的值或指针v指向的值
	e.reflectValue(v.Elem())
}
```
## 单元或集成测试
本次测试包括划分 tag、整数、浮点数、字符串、各种map、结构体、嵌套结构体、匿名结构体、解析 tag、数组。其中有一个集成测试综合测试上述需求。
### 测试代码
#### marshal_test.go
```go
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

```
#### encode_test.go
```go
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

```
#### tags_test.go
```go
package marshal

import (
	"testing"
)

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}

```
### 测试结果
```sh
[luowle@VM_0_4_centos marshal]$ go test -v
=== RUN   TestEncodeStrOfStruct
--- PASS: TestEncodeStrOfStruct (0.00s)
=== RUN   TestEncodeIntOfStruct
--- PASS: TestEncodeIntOfStruct (0.00s)
=== RUN   TestEncodeFloatOfStruct
--- PASS: TestEncodeFloatOfStruct (0.00s)
=== RUN   TestEncodeMapIntStr
--- PASS: TestEncodeMapIntStr (0.00s)
=== RUN   TestEncodeMapIntFloat
--- PASS: TestEncodeMapIntFloat (0.00s)
=== RUN   TestEncodeMapStrInt
--- PASS: TestEncodeMapStrInt (0.00s)
=== RUN   TestEncodeStructArray
--- PASS: TestEncodeStructArray (0.00s)
=== RUN   TestJsonMarshal
--- PASS: TestJsonMarshal (0.00s)
=== RUN   TestTagParsing
--- PASS: TestTagParsing (0.00s)
PASS
ok      github.com/LyleLuo/ServiceComputing/myjson/marshal      0.015s
```
## 功能测试
功能测试包括嵌套结构体、匿名结构体、tag、数组、各种基本类型
### 测试代码
```go
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

```
### 测试结果
```sh
[luowle@VM_0_4_centos myjson]$ go run main.go 
[{"ID":18342069,"Movies":{"released":2020,"color":false,"Name":"Xiyangyang","rate":4.5,"People":{"PeopleName":"Luo Weile"}},"NationName":"China"},{"ID":12345678,"Movies":{"released":2020,"color":false,"Name":"Huitailang","rate":5,"People":{"PeopleName":"Luo Luo"}},"NationName":"USA"},{"ID":98765433,"Movies":{"released":2020,"color":false,"Name":"Xiongchumo","rate":4.3,"People":{"PeopleName":"pml"}},"NationName":"UK"}]
```
## 总结
可见，各种类型的测试都能通过，可见该实现没有太大问题。本实现有较多借鉴官方库的地方，在看懂官方库的高手代码的过程中学习到了不少东西，自己也熟练的 reflect 的相关操作。最令我印象深刻的是 field 的 index 的生成与使用，本代码也参考了该实现。第一次这样通过 slice 来搜索而不是使用常用的队列，让我耳目一新。本库必要的地方都打上了注释，方便日后自己与他人观看。

