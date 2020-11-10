# 修改、改进 RxGo 包
## 设计说明文档
[设计说明文档](specification.md)
## 简单使用样例
### 代码
```go
package main

import (
	"fmt"

	"github.com/LyleLuo/ServiceComputing/RxGo/rxgo"
)

func main() {
	rxgo.Just(-1, 1, 2, 3).First().Subscribe(func(x int) {
		fmt.Println(x)
	})

	rxgo.Just("I", "I", "love", "rxgo", "rxgo", "rxgo").Distinct().Subscribe(func(x string) {
		fmt.Println(x)
	})
}
```
### 结果
```
[luowle@VM_0_4_centos RxGo]$ go run main.go 
-1
I
love
rxgo
```
[更多使用样例](main.go)