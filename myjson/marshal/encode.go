package marshal

import (
	"bytes"
	"encoding"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// EncodeState 存储v的编组之后的bytes.Buffer对象
type EncodeState struct {
	bytes.Buffer // accumulated output
	scratch      [64]byte
}

// EncoderFunc 将reflect.Value写到EncodeState的buffer函数
type EncoderFunc func(e *EncodeState, v reflect.Value)

func (e *EncodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	//继续调用relectValue方法
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *EncodeState) reflectValue(v reflect.Value) {
	//调用valueEncoder函数生成 EncoderFunc 然后执行（e,v)
	ValueEncoder(v)(e, v)
}

// ValueEncoder 通过reflect.Value的类型获取一个可返回的EncoderFunc
func ValueEncoder(v reflect.Value) EncoderFunc {
	//判断reflect.value对象是否为一个值
	if !v.IsValid() {
		return invalidValueEncoder
	}
	//继续调用typeEncoder函数
	return typeEncoder(v.Type())
}

func invalidValueEncoder(e *EncodeState, v reflect.Value) {
	e.WriteString("null")
}

var encoderCache sync.Map // map[reflect.Type]EncoderFunc

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

// 下面的函数用于返回一个错误 被解析成json的类型 没有实现
func unsupportedTypeEncoder(e *EncodeState, v reflect.Value) {
	e.error(&unsupportedTpyeError{v.Type()})
}

type unsupportedTpyeError struct {
	Type reflect.Type
}

type jsonError struct{ error }

func (e *EncodeState) error(err error) {
	panic(jsonError{err})
}

func (e *unsupportedTpyeError) Error() string {
	return "json: unsupported type: " + e.Type.String()
}

// 用于处理bool的encoderFunc
func boolEncoder(e *EncodeState, v reflect.Value) {
	if v.Bool() {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
}

// 用于处理int和uint的encoderFunc
func intEncoder(e *EncodeState, v reflect.Value) {
	b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
	e.Write(b)
}

func uintEncoder(e *EncodeState, v reflect.Value) {
	b := strconv.AppendUint(e.scratch[:0], v.Uint(), 10)
	e.Write(b)
}

type floatEncoder int // number of bits

var (
	float32Encoder = (floatEncoder(32)).encode
	float64Encoder = (floatEncoder(64)).encode
)

type unsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *unsupportedValueError) Error() string {
	return "json: unsupported value: " + e.Str
}

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

// 用于解析string的函数
func stringEncoder(e *EncodeState, v reflect.Value) {
	e.string(v.String())
}

// 将string内容写到缓冲区（加上双引号）
func (e *EncodeState) string(s string) {
	e.WriteByte('"')
	if 0 < len(s) {
		e.WriteString(s)
	}
	e.WriteByte('"')
}

// 用于解析interface类型的函数，本质上是获得接口的类型后在调用相应函数
func interfaceEncoder(e *EncodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// Elem返回接口v包含的值或指针v指向的值
	e.reflectValue(v.Elem())
}

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

// 新建一个解析器，对于结构体来说，每一种结构体都要新建一个解析方法用于缓存
func newStructEncoder(t reflect.Type) EncoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}

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

type mapEncoder struct {
	elemEnc EncoderFunc
}

func newMapEncoder(t reflect.Type) EncoderFunc {
	switch t.Key().Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		return unsupportedTypeEncoder
	}
	me := mapEncoder{typeEncoder(t.Elem())}
	return me.encode
}

func (me mapEncoder) encode(e *EncodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	e.WriteByte('{')

	// Extract and sort the keys.
	// 判断 key 是否是期待的类型，如是，则将 key 用 reflectWithString.s 表达
	keys := v.MapKeys()
	sv := make([]reflectWithString, len(keys))
	for i, v := range keys {
		sv[i].v = v
		if err := sv[i].resolve(); err != nil {
			e.error(fmt.Errorf("json: encoding error for type %q: %q", v.Type().String(), err.Error()))
		}
	}
	sort.Slice(sv, func(i, j int) bool { return sv[i].s < sv[j].s })

	// 将内容逐步写入缓冲区
	for i, kv := range sv {
		if i > 0 {
			e.WriteByte(',')
		}
		e.string(kv.s)
		e.WriteByte(':')
		// 调用 value 自身的 EncoderFunc 来写 value
		me.elemEnc(e, v.MapIndex(kv.v))
	}
	e.WriteByte('}')
}

type reflectWithString struct {
	v reflect.Value
	s string
}

// 用于对 key 进行类型检查，如果是合格类型则写到 reflectWithString.s 中
func (w *reflectWithString) resolve() error {
	if w.v.Kind() == reflect.String {
		w.s = w.v.String()
		return nil
	}
	if tm, ok := w.v.Interface().(encoding.TextMarshaler); ok {
		if w.v.Kind() == reflect.Ptr && w.v.IsNil() {
			return nil
		}
		buf, err := tm.MarshalText()
		w.s = string(buf)
		return err
	}
	switch w.v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.s = strconv.FormatInt(w.v.Int(), 10)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		w.s = strconv.FormatUint(w.v.Uint(), 10)
		return nil
	}
	panic("unexpected map key type")
}

type arrayEncoder struct {
	elemEnc EncoderFunc
}

// array 的实现最为简单，只需要把中括号写出来，然后边调用数组的encoderFunc一边输出逗号
func (ae arrayEncoder) encode(e *EncodeState, v reflect.Value) {
	e.WriteByte('[')
	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteByte(',')
		}
		ae.elemEnc(e, v.Index(i))
	}
	e.WriteByte(']')
}

func newArrayEncoder(t reflect.Type) EncoderFunc {
	enc := arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

// 用于缓存 field
var fieldCache sync.Map // map[reflect.Type]structFields

// 缓存 field 函数从而避免重复操作
func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}

func typeFields(t reflect.Type) structFields {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				// 如果是匿名且未导出的非结构体就跳过，匿名结构体不跳过的原因是有可能里面有可导出type
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
				// 对tag进行解析
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				// 将当前元素记录，len(index)表示所在层数
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
					}
					field.nameBytes = []byte(field.name)

					field.nameNonEsc = `"` + field.name + `":`

					fields = append(fields, field)
					continue
				}

				// 记录下匿名结构体
				next = append(next, field{name: ft.Name(), index: index, typ: ft})
			}
		}
	}

	// 把刚刚解析的结果都记录到 structFields 里
	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}
	nameIndex := make(map[string]int, len(fields))
	for i, field := range fields {
		nameIndex[field.name] = i
	}
	return structFields{fields, nameIndex}
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

// 判断当前tag是否为合法tag
func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

// 判断当前值是否为空值
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
