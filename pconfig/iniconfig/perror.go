package iniconfig

import (
	"fmt"
	"os"
)

// Perror 用于输出错误信息并推出程序的函数，参数err是错误
func Perror(err error) {
	if err != nil {
		if _, err2 := fmt.Fprintln(os.Stderr, err); err2 != nil {
			panic(err2)
		}
		os.Exit(1)
	}
}
