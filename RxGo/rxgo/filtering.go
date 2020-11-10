package rxgo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

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

// First 返回第一位的元素，即Element[0]
func (parent *Observable) First() (o *Observable) {
	return parent.ElementAt(0)
}

// IgnoreElements 忽略所有元素
func (parent *Observable) IgnoreElements() (o *Observable) {
	o = parent.newTransformObservable("IgnoreElements")
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			return
		}}
	return o
}

// Last 发射最后一项
func (parent *Observable) Last() (o *Observable) {
	return parent.Takelast(1)
}

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
