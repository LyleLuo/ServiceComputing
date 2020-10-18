# specification
## 设计说明
### 核心功能 Watch
Watch 用来监听文件改变的函数，filename 是文件路径，listener 是接口，返回一个储存文件信息的 map 和错误类型 error。Watch 里调用 listen 对文件修改情况进行轮询，如果 listen 轮询结束将修改后的配置以 map[string]map[string]string 返回。
```go
// Watch 用来监听文件改变的函数，filename 是文件路径，listener 是接口，返回一个储存文件信息的 map 和错误类型 error
func Watch(filename string, listener Listener) (map[string]map[string]string, error) {
	listener.listen(filename)
	c, err := InitConfig(filename)
	if err != nil {
		Perror(err)
	}
	return c.Conflist, err
}
```
Watch返回值的使用
```ini
[paths]
data = /home/git/grafana
```
第一个 string 对应上面的 paths，第二个 string 对应上面的 data，第三个 string 对应着 /home/git/grafana
```go
map[string]map[string]string
```

### 配置结构体
```go
// Config 储存配置文件的信息，包括文件名和文件里的键值对
type Config struct {
	filepath string
	Conflist map[string]map[string]string
}
```

### 接口的定义
```go
// Listener 接口，所有实现了接口的函数的类型都可以转换为接口调用对应的函数
type Listener interface {
	listen(inifile string)
}

// ListenFunc 是一种函数类型，其格式为func(string)
type ListenFunc func(string)

func (lf ListenFunc) listen(inifile string) {
	lf(inifile)
}
```

### 接口的使用
所有满足签名的函数、方法都可以作为参数。
下面是满足签名的函数
```go
// Listen 参数inifile是要监听改变的文件，该函数使用轮询的方式，直到被监听的文件发生了改变
func Listen(inifile string) {
	originConf, err := InitConfig(inifile)
	if err != nil {
		Perror(err)
	}
	for {
		time.Sleep(time.Duration(1) * time.Second)
		newConf, err := InitConfig(inifile)
		if err != nil {
			Perror(err)
		}
		if !reflect.DeepEqual(originConf, newConf) {
			break
		}
	}
}
```
所有实现 Listener 接口的数据类型都可作为参数。下面是 Config 结构体实现 listen
```go
func (c *Config) listen(inifile string) {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		newConf, err := InitConfig(inifile)
		if err != nil {
			Perror(err)
		}
		if !reflect.DeepEqual(c, newConf) {
			break
		}
	}
}
```
上面两者都可以使用 Listener.listen(ininame string) 调用。例子如下
#### 使用结构体式接口
```go
func main() {
    // InitConfig 返回一个*Config和错误
	c1, err := iniconfig.InitConfig("test.ini")
	m1 := c1.Conflist
	if err != nil {
		iniconfig.Perror(err)
	}
    fmt.Println(m1)
    
    // c1 *Config 实现了 listen 方法，因此可以传入 Watch
	m2, err := iniconfig.Watch("test.ini", c1)
	if err != nil {
		iniconfig.Perror(err)
	}
	fmt.Println(m2)
}
```

#### 使用函数式接口
```go
func main() {
	var f iniconfig.ListenFunc = Listen
	m2, err := iniconfig.Watch("test.ini", f)
	if err != nil {
		iniconfig.Perror(err)
	}
	fmt.Println(m2)

}
```

### 自定义错误
由于如配置文件不存在，无法打开等错误在 bufio 中已经提供，此处不再多此一举。

Config 的 GetValue 有可能不存在 key 从而无法获取值，于是自定义了一个错误
```go
// ErrKeyNotExist 错误：不存在该键值对
var ErrKeyNotExist = errors.New("The section don't include this key-value pair")
```
当尝试用一些不存在的 key 试图获取 value 时会返回此错误
```go
// GetValue 通过参数section和name获得值，返回string是value，error是错误类型
func (c *Config) GetValue(section, name string) (string, error) {
	if _, ok := c.Conflist[section][name]; ok {
		return c.Conflist[section][name], nil
	}
	return "", ErrKeyNotExist
}
```

### init 函数
使得 Unix 系统默认采用 # 作为注释行，Windows 系统默认采用 ; 作为注释行
```go
// 在读取文件时根据 commentChar 判断注释符号
var commentChar rune

func init() {
	if runtime.GOOS == "windows" {
		commentChar = ';'
	} else {
		commentChar = '#'
	}
}
```

## 单元或集成测试
### 测试代码
config_test.go
```go
package iniconfig

import "testing"

func TestInitConfig(t *testing.T) {
	c, err := InitConfig("../test.ini")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if c.Conflist[""]["app_mode"] != "development" {
		t.Errorf("Not equal: expected: development, actual: %s", c.Conflist[""]["app_mode"])
	}
	if c.Conflist["paths"]["data"] != "/home/git/grafana" {
		t.Errorf("Not equal: expected: /home/git/grafana, actual: %s", c.Conflist[""]["app_mode"])
	}
	if c.Conflist["server"]["protocol"] != "http" {
		t.Errorf("Not equal: expected: http, actual: %s", c.Conflist[""]["app_mode"])
	}
	if c.Conflist["server"]["http_port"] != "9999" {
		t.Errorf("Not equal: expected: 9999, actual: %s", c.Conflist[""]["app_mode"])
	}
	if c.Conflist["server"]["enforce_domain"] != "true" {
		t.Errorf("Not equal: expected: true, actual: %s", c.Conflist[""]["app_mode"])
	}
}

func TestGetValue(t *testing.T) {
	c, err := InitConfig("../test.ini")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if value, err := c.GetValue("server", "protocol"); value != "http" || err != nil {
		t.Errorf("Not equal: expected: http, actual: %s", c.Conflist[""]["app_mode"])
	}
	if value, err := c.GetValue("server", "haha"); value != "" || err != ErrKeyNotExist {
		t.Errorf("Return error failed")
	}
}
```
perror_test.go
```go
package iniconfig

import "testing"

func TestPerror(t *testing.T) {
	Perror(nil)
}
```
watch_listen_test.go
```go
package iniconfig

import "testing"

// 运行此测试的过程中把http_port = 9999改为http_port = 8888
func TestWatch(t *testing.T) {
	var f ListenFunc = Listen
	m, err := Watch("../test.ini", f)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if m[""]["app_mode"] != "development" {
		t.Errorf("Not equal: expected: development, actual: %s", m[""]["app_mode"])
	}
	if m["paths"]["data"] != "/home/git/grafana" {
		t.Errorf("Not equal: expected: /home/git/grafana, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["protocol"] != "http" {
		t.Errorf("Not equal: expected: http, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["http_port"] != "8888" {
		t.Errorf("Not equal: expected: 8888, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["enforce_domain"] != "true" {
		t.Errorf("Not equal: expected: true, actual: %s", m[""]["app_mode"])
	}
}
```

### 测试结果
由于本处有对 Watch 进行简单的测试，go test 后会进行轮询，需要手动在运行此测试的过程中把http_port = 9999 改为 http_port = 8888。如果修改成功，测试则得以通过。

（可以看出来，我跑去改ini文件的时间大概是五秒）
```sh
[luowle@VM_0_4_centos iniconfig]$ go test
PASS
ok      github.com/LyleLuo/ServiceComputing/pconfig/iniconfig   5.007s
```

## 功能测试
### 正确性测试
测试主函数如下所示
```go
package main

import (
	"fmt"

	"github.com/LyleLuo/ServiceComputing/pconfig/iniconfig"
)

func main() {
	c1, err := iniconfig.InitConfig("test.ini")
	m1 := c1.Conflist
	if err != nil {
		iniconfig.Perror(err)
	}
	fmt.Println(m1)
	m2, err := iniconfig.Watch("test.ini", c1)
	// var f iniconfig.ListenFunc = iniconfig.Listen
	// m2, err := iniconfig.Watch("test.ini", f)
	if err != nil {
		iniconfig.Perror(err)
	}
	fmt.Println(m2)

}
```
开始执行
```sh
[luowle@VM_0_4_centos pconfig]$ go run main.go 
map[:map[app_mode:development] paths:map[data:/home/git/grafana] server:map[enforce_domain:true http_port:9999 protocol:http]]

```
原来的文件
```ini
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = http

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
```
修改为
```ini
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = https

# The http port  to use
https_port = 443

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
```
运行结果如下
```sh
[luowle@VM_0_4_centos pconfig]$ go run main.go 
map[:map[app_mode:development] paths:map[data:/home/git/grafana] server:map[enforce_domain:true http_port:9999 protocol:http]]
map[:map[app_mode:development] paths:map[data:/home/git/grafana] server:map[enforce_domain:true https_port:443 protocol:https]]
```
可见实现了所需的功能，Watch 返回了一个新的配置信息。

### 错误性测试
修该主函数，给一个不存在的文件路径看看是否报错。
```go
func main() {
    // 只修改了下面这行，给了个错误路径
	c1, err := iniconfig.InitConfig("hhhhhhh")
	m1 := c1.Conflist

```
如期产生错误信息并退出
```sh
[luowle@VM_0_4_centos pconfig]$ go run main.go 
open hhhhhhh: no such file or directory
exit status 1
```

### 功能测试总结
可见，当路径正确时，修改文件后 Watch 返回了一个新的配置信息。而当路径错误时，也可以返回正确的错误。该包的功能较为完备。

