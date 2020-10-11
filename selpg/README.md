# 中山大学服务计算第三次作业
## 实验准备
安装pflag
```sh
go get github.com/spf13/pflag
```
## 设计说明
该项目主要功能由pargs()和pinput()两个函数提供，分别用于处理参数和处理输入文件
### 自定义类型
设计了一个关于输入参数的结构体，各参数作用可以顾名思义。
```go
type Page struct {
	StartPage  int
	EndPage    int
	InFileName string
	PageLen    int
	PageType   bool
	PrintDest  string
}
```
定义了一些错误
```go
var errPageNotPos = errors.New("start page and end page should be positive")
var errPageStartGreater = errors.New("end page should be greater then start page")
var errTooManyArgs = errors.New("there are too many arguments")
var errFile = errors.New("cannot open your input file")
var errPipe = errors.New("cannot open pipe")
```

### pargs函数
使用pflag解析参数，如果出现逻辑上的错误则返回上面所定义的错误
```go
	pflag.IntVarP(&page.StartPage, "start-page", "s", 0, "start page")
	pflag.IntVarP(&page.EndPage, "end-page", "e", 0, "end page")
	pflag.IntVarP(&page.PageLen, "page-len", "l", 72, "page length")
	pflag.BoolVarP(&page.PageType, "page-type", "f", false, "page type")
	pflag.StringVarP(&page.PrintDest, "print-dest", "d", "", "printer destination")
	pflag.Usage = func() {
		usage()
		pflag.PrintDefaults()
	}
	pflag.Parse()
```

### pinput函数
#### 设置在哪里读取数据
```go
if len(page.InFileName) > 0 {
		fIn, err := os.Open(page.InFileName)
		if err != nil {
			return errFile
		}
		defer fIn.Close()
		read = bufio.NewReader(fIn)
	} else {
		read = bufio.NewReader(os.Stdin)
	}

```
#### 设置数据写到哪里
```go
if len(page.PrintDest) > 0 {
		cmd := exec.Command("lp", "-d"+page.PrintDest)
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			return errPipe
		}
		defer stdinPipe.Close()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Start(); err != nil {
			return err
		}
		write = bufio.NewWriter(stdinPipe)
	} else if outputBuf != nil {
		write = bufio.NewWriter(outputBuf)
	} else {
		write = bufio.NewWriter(os.Stdout)
	}
	defer write.Flush()
```

#### 根据PageType来选择如何读写
```go
if page.PageType {
		for {
			ch, err := read.ReadByte()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			if page.StartPage <= pages && pages <= page.EndPage {
				if err := write.WriteByte(ch); err != nil {
					return err
				}
			}
			if ch == '\f' {
				pages++
			}
		}
	} else {
		for {
			// 这次可以使用封装程度更高的scan
			cstream, err := read.ReadBytes('\n')
			if err == io.EOF {
				if page.StartPage <= pages && pages <= page.EndPage {
					if _, err := write.Write(cstream); err != nil {
						return err
					}
				}
				break
			} else if err != nil {
				return err
			}

			if lines > page.PageLen {
				pages++
				lines = 1
			}
			if page.StartPage <= pages && pages <= page.EndPage {
				if _, err := write.Write(cstream); err != nil {
					return err
				}
			}
			lines++
		}
	}
```

### 为进行测试而进行的考虑
由于该项目是直接输出到stdout和stderr中的，因此为了进行测试，需要进行额外的设计。

调用时给程序一个缓冲区
```go
	buf := new(bytes.Buffer)
	Selpg(page, buf)
```

Selpg调用pinput同样需要传入这个缓冲区，之后判断如果存在缓冲区则把输出写到缓冲区。
```go
} else if outputBuf != nil {
    write = bufio.NewWriter(outputBuf)
```

## 测试
### 单元测试
设置一个用于测试的结构体和一个缓冲区传入Selpg()，将缓冲区的输出和期望输出作比对。该测试主要测试pinput的逻辑正确。
```go
package myselpg

import (
	"bytes"
	"testing"
)

func TestSelpg(t *testing.T) {
	page := new(Page)
	page.StartPage = 4
	page.EndPage = 5
	page.InFileName = "../test/testfile"
	page.PageLen = 2
	buf := new(bytes.Buffer)
	Selpg(page, buf)
	output := buf.String()
	exceptOutput := "\t\"io\"\n\t\"os\"\n\t\"os/exec\"\n\n"
	if output != exceptOutput {
		t.Errorf("Excepted output is:\n%s, but you output is:\n%s", exceptOutput, output)
	}
}
```
测试结果如下所示
```sh
[luowle@VM_0_4_centos selpg]$ cd myselpg/
[luowle@VM_0_4_centos myselpg]$ go test
PASS
ok      github.com/LyleLuo/ServiceComputing/selpg/myselpg       0.005s
```
### 集成测试
此处自己编写了一个脚本进行了自动化测试。思路是将所有输出重定向到一个output文件里，再使用diff和期望输出作比较。脚本如下
```sh
selpg -s2 -e3 testfile >output1
diff output1 ans1

selpg -s10 -e20 -l10 testfile >output2
diff output2 ans2

selpg -s10 -e20 -l10 <testfile >output2
diff output2 ans2

selpg -s10 -e20 -l10 testfile | cat >output2
diff output2 ans2

cat testfile | selpg -s10 -e20 -l10 >output2
diff output2 ans2

selpg -s10 -e20 -l10 null 2>output3
diff output3 ans3

selpg -s10 -e20 -l10 testfile >output2 2>output3
diff output2 ans2

selpg -s10 -e20 -l10 testfile >output2 2>/dev/null
diff output2 ans2

selpg -s10 -e20 -l10 testfile >/dev/null
diff output2 ans2

selpg -s3 -e4 -f testfile >output4
diff output4 ans4
```
如果所有结果都和预期输出的一样，那么执行该脚本则不会有任何输出。执行结果如下所示
```sh
[luowle@VM_0_4_centos selpg]$ cd test/
[luowle@VM_0_4_centos test]$ sh test.sh 
[luowle@VM_0_4_centos test]$ 
```

### 功能测试
1. 
```
[luowle@VM_0_4_centos test]$ selpg -s1 -e1 testfile 
package myselpg

import (
        "bufio"
        "errors"
        "fmt"
        "io"
        "os"
        "os/exec"

        "github.com/spf13/pflag"
)

// Page 是selpg的命令行参数结构体
var errPageStartGreater = errors.New("end page should be greater then start page")
var errTooManyArgs = errors.New("there are too many arguments")
var errFile = errors.New("cannot open your input file")
var errPipe = errors.New("cannot open pipe")

func perror(err error) {
        if err != nil {
                if _, err2 := fmt.Fprintln(os.Stderr, err); err2 != nil {
                        panic(err2)
                }
                os.Exit(1)
        }
}

func usage() {

        progName := os.Args[0]
        fmt.Fprintf(os.Stderr, "USAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progName)
}

func pargs(page *Page) error {
        pflag.IntVarP(&page.StartPage, "start-page", "s", 0, "start page")
        pflag.IntVarP(&page.EndPage, "end-page", "e", 0, "end page")
        pflag.IntVarP(&page.PageLen, "page-len", "l", 72, "page length")
        pflag.BoolVarP(&page.PageType, "page-type", "f", false, "page type")
        pflag.StringVarP(&page.PrintDest, "print-dest", "d", "", "printer destination")
        pflag.Usage = func() {
                usage()
                pflag.PrintDefaults()
        }
        pflag.Parse()
        if page.StartPage <= 0 || page.EndPage <= 0 {
                pflag.Usage()
                return errPageNotPos
        }

        if page.StartPage > page.EndPage {
                return errPageStartGreater
        }

        if pflag.NArg() == 1 {
                page.InFileName = pflag.Arg(0)
        } else if pflag.NArg() > 1 {
                return errTooManyArgs
        } else {
                page.InFileName = ""
        }

        return nil
[luowle@VM_0_4_centos test]$ 

```

2. 
```
[luowle@VM_0_4_centos test]$ selpg -s1 -e1 <testfile 
package myselpg

import (
        "bufio"
        "errors"
        "fmt"
        "io"
        "os"
        "os/exec"

        "github.com/spf13/pflag"
)
...太长了不重复展示，输出和第一个一样
```

3. 
```
[luowle@VM_0_4_centos test]$ cat testfile | selpg -s10 -e20
```
结果太长了，补上>func_test3重定向结果到func_test3上了

4. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile >func_test4
```
重定向结果到func_test4

5. 
```
selpg -s10 -e20 testfile 2>error_file_5
```
结果太长了，补上>func_test5重定向结果到func_test5上了。由于没有错误，error_file_5是空的

6. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile >func_test6 2>error_file_6
```
结果在func_test6和error_file_6

7. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile >func_test7 2>/dev/null
```
结果在func_test7

8. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile >/dev/null
```
输出结果被丢弃

9. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile | cat 
```
结果太长了，补上>func_test9重定向结果到func_test9上了

10. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile 2>error_file_10
```
结果太长了，补上>func_test10重定向结果到func_test10上了

11. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 -l3 testfile 
var errPipe = errors.New("cannot open pipe")

func perror(err error) {
        if err != nil {
                if _, err2 := fmt.Fprintln(os.Stderr, err); err2 != nil {
                        panic(err2)
                }
                os.Exit(1)
        }
}

func usage() {

        progName := os.Args[0]
        fmt.Fprintf(os.Stderr, "USAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progName)
}

func pargs(page *Page) error {
        pflag.IntVarP(&page.StartPage, "start-page", "s", 0, "start page")
        pflag.IntVarP(&page.EndPage, "end-page", "e", 0, "end page")
        pflag.IntVarP(&page.PageLen, "page-len", "l", 72, "page length")
        pflag.BoolVarP(&page.PageType, "page-type", "f", false, "page type")
        pflag.StringVarP(&page.PrintDest, "print-dest", "d", "", "printer destination")
        pflag.Usage = func() {
                usage()
                pflag.PrintDefaults()
        }
        pflag.Parse()
        if page.StartPage <= 0 || page.EndPage <= 0 {
                pflag.Usage()
                return errPageNotPos
        }

        if page.StartPage > page.EndPage {
```

12. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 -f testfile >func_test12
```
结果太长了，补上>func_test12重定向结果到func_test12上了

13. 由于lp命令和打印机功能在本机都不可行，修改源代码使其调用其他命令

原来的代码
```go
	if len(page.PrintDest) > 0 {
		cmd := exec.Command("lp", "-d"+page.PrintDest)
```

现在的代码
```go
	if len(page.PrintDest) > 0 {
		cmd := exec.Command("cat", "-n")
```
修改后，只要传入-d参数就可以使用管道传进cat命令
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 -dcat testfile >func_test13
```
结果太长了，补上>func_test13重定向结果到func_test13上了

14. 
```
[luowle@VM_0_4_centos test]$ selpg -s10 -e20 testfile >func_test14  2>error_file_14 &
[1] 13571
[luowle@VM_0_4_centos test]$ ps
  PID TTY          TIME CMD
 9047 pts/5    00:00:00 bash
11621 pts/5    00:00:00 bash
13609 pts/5    00:00:00 ps
[1]+  Done                    selpg -s10 -e20 testfile > func_test14 2> error_file_14
```
结果在func_test14和error_file_14

其中，测试结果3 4 5 6 7 9 10 14是**一样**的。 
所有重定向的记录都在**test**文件夹中。  
所有的测试均符合预期，程序基本是正确的。
