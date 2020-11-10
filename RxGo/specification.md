# specification

## 设计说明

### 文件
#### 新增的文件
- filtering.go -- filtering 操作所在的文件
- filtering_test.go -- filtering 操作对应的测试文件

#### 修改的文件
- rxgo.go -- 修改 Observable 结构体以满足设计需求
- transforms.go -- 修改 op 函数以满足设计需求

### 设计思路
需要实现的操作如下

Filtering Observables

Operators that selectively emit items from a source Observable.

- Debounce — only emit an item from an Observable if a particular timespan has passed without it emitting another item
- Distinct — suppress duplicate items emitted by an Observable
- ElementAt — emit only item n emitted by an Observable
- Filter — emit only those items from an Observable that pass a predicate test
- First — emit only the first item, or the first item that meets a condition, from an Observable
- IgnoreElements — do not emit any items from an Observable but mirror its termination notification
- Last — emit only the last item emitted by an Observable
- Sample — emit the most recent item emitted by an Observable within periodic time intervals
- Skip — suppress the first n items emitted by an Observable
- SkipLast — suppress the last n items emitted by an Observable
- Take — emit only the first n items emitted by an Observable
- TakeLast — emit only the last n items emitted by an Observable

其中，Filter 操作已在库中实现。因此我只需要参考 Filter 实现其他操作。

### 源代码分析

- 函数最后的操作需要传入函数到 Subscribe 中
- 流的发送也是在调用 Subscribe 才开始的
- Subscribe 将各个 Observable 连接起来
- Subscribe 在最后一个流输出中使用 range 遍历 outflow，并用传入的函数进行操作

```go
func (o *Observable) Subscribe(ob interface{}) {
	......
	// 将各个 Observable 连接起来
	o.connect(ctx)
	if ctxok {
		oc.OnConnected()
	}

	// 获得最后一个 observable 的outflow 作为本函数的输入流
	po := o
	for ; po.next != nil; po = po.next {
	}

	in := po.outflow
	o.mu.Unlock()
	// 使用 fv（传入的函数）将最后一个输出流进行操作
	for x := range in {
		if observer != nil {
			if e, ok := x.(error); ok {
				observer.OnError(e)

			} else {
				observer.OnNext(x)
			}
		} else {
			if _, ok := x.(error); ok {
				// skip error
			} else {
				params := []reflect.Value{reflect.ValueOf(x)}
				fv.Call(params)
			}
		}
	}
	if observer != nil {
		observer.OnCompleted()
	}
}
```

其中，connect 函数操作如下所示。根据源代码很容易看出来，该函数时从 Observable 根开始，为每个 Observable 创建输出的 channel，然后调用 Observable.operator.op 进行操作。
```go
// connect all Observable form the first one.
func (o *Observable) connect(ctx context.Context) {
	for po := o.root; po != nil; po = po.next {
		po.outflow = make(chan interface{}, po.buf_len)
		po.operator.op(ctx, po)
	}
}
```

以 Just 为例，其 op 函数如下所示。
```go
func (sop sourceOperater) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	out := o.outflow
	//fmt.Println(o.name, "source out chan ", out)

	// Scheduler
	go func() {
		for end := false; !end; { // made panic op re-enter
			end = sop.opFunc(ctx, o, out)
		}
		o.closeFlow(out)
	}()
}
```
可以看出，sourceOperater 的 op 会直接调用其 opFunc。而 Just 的 opFunc 如下所示。该函数发送完所有的元素就会返回 end = true，然后上面的 op 就会关闭流通道，使得接收者可以跳出 for 循环停止监听。
```go
var rangeSource = sourceOperater{func(ctx context.Context, o *Observable, out chan interface{}) (end bool) {
	fv := reflect.ValueOf(o.flip)
	params := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(out)}
	// 调用用于发射的函数
	fv.Call(params)
	// 返回结束标志
	return true
}}

// Just creates an Observable with the provided item(s).
func Just(items ...interface{}) *Observable {
	o := newGeneratorObservable("Just")
	// 将所有元素发送的函数，不过并不在这里使用，而是在 op 中使用
	o.flip = func(ctx context.Context, out chan interface{}) {
		for _, item := range items {
			if b := o.sendToFlow(ctx, item, out); b {
				return
			}
		}
	}
	o.operator = justSource
	return o
}

var justSource = rangeSource
```

如果 Just 之后没有其他操作了，connect 的循环只会执行一次，然后 sourceOperater 所产生的 outflow 可以直接被 Subscribe 调用所传进的参数函数使用。但如果 Just 之后还有其他操作，假设 transOperater 是其下一个操作，则 connect 将 上一个输出 chan 作为 transOperater.op 的输入 chan，然后 transOperater.op 会产生一个 go 程对输入流进行轮询，每当有输入就调用 Observable.operation，输入通道关闭后则停止轮询，go 程退出。transOperater.op 代码如下所示，可以看出，该 op 函数可以基本满足完成 filtering 操作的需求。
```go
func (tsop transOperater) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup
	// 产生一个 go 程对输入流进行轮询
	go func() {
		end := false
		for x := range in {
			if end {
				continue
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			// scheduler
			switch threading := o.threading; threading {
				// 每当有输入就调用 Observable.operation
				case ThreadingDefault:
					if tsop.opFunc(ctx, o, xv, out) {
						end = true
					}
				case ThreadingIO:
					fallthrough
				case ThreadingComputing:
					wg.Add(1)
					go func() {
						defer wg.Done()
						if tsop.opFunc(ctx, o, xv, out) {
							end = true
						}
					}()
				default:
			}
		}

		wg.Wait() //waiting all go-routines completed
		o.closeFlow(out)
	}()
}
```

以 Filter 为例，Filter 首先创建一个 Observable，然后把 filterOperater 赋给 Observable.operator。使得上面 transOperater 可以使用。
```go
// Filter `func(x anytype) bool` filters items in the original Observable and returns
// a new Observable with the filtered items.
func (parent *Observable) Filter(f interface{}) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(f)
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("filter")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = filterOperater
	return o
}
```

filterOperater 的实现就是靠下面函数。一句话概括就是如果满足之前传入的筛选函数（o.flip）的条件，就使用 o.sendToFlow 把数据发送到输出通道。
```go
var filterOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}
	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
	}

	return
}}
```

### 核心的设计与实现

上面经过源代码分析可知，只需要将目标功能的函数并实现其 Observable.operation 即可。

#### Debounce
设计思路：使用 tempCount 记录当前是第几个被发射的，count 是当前被发射总数。如果休眠一段时间后 tempCount 与 count 相等，则说明这段时间没有其他元素被发射，该元素可以被发射。代码如下所示

注意：该处使用ThreadingComputing，即 go 程并发操作，因此count 的操作要使用互斥锁。

```go
// Debounce  仅在经过特定时间跨度时才从Observable发出一项，而不发出另一项
func (parent *Observable) Debounce(stopTime time.Duration) (o *Observable) {
	o = parent.newTransformObservable("Debounce")
	o.threading = ThreadingComputing
	count := 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			var mutex sync.Mutex
			mutex.Lock()
			count++
			tempCount := count
			mutex.Unlock()
			time.Sleep(stopTime)
			mutex.Lock()
			if tempCount == count {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			mutex.Unlock()
			return
		}}
	return o
}
```

#### Distinct
设计思路：使用一个 map 记录当前已出现过的 item。如果 item 没有出现过，则加入 map 中标记并将其发射，否则舍弃该 item。代码如下所示
```go
// Distinct 消除重复项
func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newTransformObservable("Distinct")
	distinctMap := map[string]bool{}
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			itemStr := fmt.Sprintf("%v", item)
			if _, ok := distinctMap[itemStr]; !ok {
				distinctMap[itemStr] = true
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			return
		}}
	return o
}
```

#### ElementAt 
设计思路：使用 count 记录当前已到达的 item。如果 count 等于所传入的 id，则将其发射，否则舍弃。
```go
// ElementAt 返回位于第id位的元素，从0开始
func (parent *Observable) ElementAt(id int) (o *Observable) {
	o = parent.newTransformObservable("ElementAt")
	count := 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if count == id {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			count++
			return
		}}
	return o
}
```

#### Filter
原本包中已实现，该实现代码是本次扩展的重要参考

#### First 
设计思路：直接复用上面的 ElementAt 即可
```go
// First 返回第一位的元素，即Element[0]
func (parent *Observable) First() (o *Observable) {
	return parent.ElementAt(0)
}
```

#### IgnoreElements
设计思路：忽略所有的值，因此 transOperater 直接返回即可
```go
// IgnoreElements 忽略所有元素
func (parent *Observable) IgnoreElements() (o *Observable) {
	o = parent.newTransformObservable("IgnoreElements")
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			return
		}}
	return o
}
```

#### Last 
设计思路：直接复用 Takelast，传入参数 1
```go
// Last 发射最后一项
func (parent *Observable) Last() (o *Observable) {
	return parent.Takelast(1)
}
```
#### Sample
设计思路：
1. 当接收到第一个的 item 时，创建一个 go 程开始循环采样计时休眠。每接收到 item 就让计数变量 count++，当每次循环采样计时休眠结束后，则对当前的最后接收到(第 count - 1) 的 item 发射出去，如果没有接收到就不发射。
2. 这些到来的 item 有并发需求，因此需要使用 ThreadingComputing 模式。
3. item 接收完成后需要将之前用来采样的 go 程关掉，此处使用 context。因此我在 transOperator.op 中添加了如下代码，使用 WithCancel 使得用于采样的子程可以随父程关闭

```go
func (tsop transOperater) op(ctx context.Context, o *Observable) {
	......
	cancelCtx, cancel := context.WithCancel(ctx)
	go func() {
		end := false
		for x := range in {
		......
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if tsop.opFunc(cancelCtx, o, xv, out) {
						end = true
					}
				}()
			default:
		......
		cancel()
	}()
}
```

注意：使用并发的 go 程需要加互斥锁
```go
// Sample 定期发射数据
func (parent *Observable) Sample(stopTime time.Duration) (o *Observable) {
	o = parent.newTransformObservable("Sample")
	o.threading = ThreadingComputing
	temp := make([]reflect.Value, 1024)
	var istimer = false
	count := 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			var mutex sync.Mutex
			mutex.Lock()
			defer mutex.Unlock()
			temp[count] = item
			count++
			if !istimer {
				istimer = true
				go func() {
					for {
						time.Sleep(stopTime)
						select {
						case <-ctx.Done():
							return
						default:
							if count != 0 && o.sendToFlow(ctx, temp[count-1].Interface(), out) {
								return
							}
							count = 0
						}
					}
				}()
			}
			time.Sleep(stopTime)
			return
		}}
	return o
}
```
#### Skip
设计思路：使用 count 记录当前已到达的 item，当 count 大于要跳过的 n 时再把 item 发送
```go
// Skip 跳过前n项再发送
func (parent *Observable) Skip(n int) (o *Observable) {
	o = parent.newTransformObservable("Skip")
	var count = 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if count >= n {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			count++
			return
		}}
	return o
}

```
#### Skiplast
设计思路：使用一个定长为 n 的队列，里面是不允许发射的 item。一开始把前 n 个到达的 item 放入队列中，再有 item 到达时，令队首出队并发射，然后把刚到达的 item 放入队尾，直到没有 item 到达为止。
```go
// Skiplast 跳过后n项后再发送
func (parent *Observable) Skiplast(n int) (o *Observable) {
	o = parent.newTransformObservable("Skiplast")
	var skipedItem []reflect.Value
	var count = 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if count == n {
				end = o.sendToFlow(ctx, skipedItem[0].Interface(), out)
				skipedItem = skipedItem[1:]
			} else {
				count++
			}
			skipedItem = append(skipedItem, item)
			return
		}}
	return o
}
```
#### Take
设计思路：使用 count 记录当前已到达的 item，当 count 小于要发送的前 n 时把 item 发送
```go
// Take 只发射前n项
func (parent *Observable) Take(n int) (o *Observable) {
	o = parent.newTransformObservable("Take")
	count := 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if count < n {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			count++
			return
		}}
	return o
}
```
#### Takelast
设计思路：像 Skiplast 一样使用一个队列保存后 n 个 item，不过 Skiplast 队列中的 item 是要舍弃掉的，而在 Takelast 中这个队列是要被发射的。

设计要点：由于 Observable.operator 是无法判断什么时候是最后到达的 item，而在 transOperater.op 则知道什么时候 inflow 的 channel 关闭。因此我们的目标是在 transOperater.op 中发射这个队列。但是现有的结构无法完成这个事情，因此需要增加一个缓冲区，使得 Observable.operator 把队列放入此处。
```go
type Observable struct {
	......
	buffer   interface{} // 新增缓冲区
	......
}

```
然后在 transOperater.op 中使用该缓冲区。如果该缓冲区非空，就在发射里面的元素
```go
func (tsop transOperater) op(ctx context.Context, o *Observable) {
		......
		if o.buffer != nil {
			buffer, ok := o.buffer.([]reflect.Value)
			if !ok {
				panic("Can't read the buffer used by recoeding last")
			}
			for _, v := range buffer {
				o.sendToFlow(ctx, v.Interface(), out)
			}
		}
		.....
}
```
Takelast 的代码与 Skiplast相似
```go
// Takelast 只发射后n项
func (parent *Observable) Takelast(n int) (o *Observable) {
	o = parent.newTransformObservable("Takelast")
	var selected []reflect.Value
	var selectedCount = 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if selectedCount == n {
				selected = selected[1:]
			} else {
				selectedCount++
			}
			selected = append(selected, item)
			o.buffer = selected
			return
		}}
	return o
}
```


## 单元或集成测试

### 测试代码
```go
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
```
### 测试结果
```sh
[luowle@VM_0_4_centos rxgo]$ go test -v filtering_test.go 
=== RUN   TestDebounce
--- PASS: TestDebounce (0.13s)
=== RUN   TestDistinct
--- PASS: TestDistinct (0.00s)
=== RUN   TestElementAt
--- PASS: TestElementAt (0.00s)
=== RUN   TestFirst
--- PASS: TestFirst (0.00s)
=== RUN   TestIgnoreElements
--- PASS: TestIgnoreElements (0.00s)
=== RUN   TestLast
--- PASS: TestLast (0.00s)
=== RUN   TestSample
--- PASS: TestSample (1.11s)
=== RUN   TestSkip
--- PASS: TestSkip (0.00s)
=== RUN   TestSkiplast
[4/4]0xc000016280
--- PASS: TestSkiplast (0.00s)
=== RUN   TestTake
--- PASS: TestTake (0.00s)
=== RUN   TestTakelast
--- PASS: TestTakelast (0.00s)
PASS
ok      command-line-arguments  (cached)
```
## 功能测试
### 测试代码
```go
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

```
### 测试结果
测试结果符合预期
```sh
[luowle@VM_0_4_centos RxGo]$ go run main.go 
[3 9]
[rxgo]
rxgo
-1
[]
lwl
[2 4 6 8 9]
[4 5 6 7 8 9]
[0 1 2 3]
[i love you TA]
[7 8 9]
```
## 总结
可见，各种类型的测试都能通过，功能测试的结果也符合预期，可见该实现没有太大问题。本实现主要借鉴原本就已经实现的 Filter 功能，在看懂潘老师的代码的过程中学习到了不少东西，自己也熟练了 context 、锁、go 程、 chan的相关操作。
