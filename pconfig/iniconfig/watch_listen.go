package iniconfig

import (
	"reflect"
	"time"
)

// Watch 用来监听文件改变的函数，filename 是文件路径，listener 是接口，返回一个储存文件信息的 map 和错误类型 error
func Watch(filename string, listener Listener) (map[string]map[string]string, error) {
	listener.listen(filename)
	c, err := InitConfig(filename)
	if err != nil {
		Perror(err)
	}
	return c.Conflist, err
}

// Listener 接口，所有实现了接口的函数的类型都可以转换为接口调用对应的函数
type Listener interface {
	listen(inifile string)
}

// ListenFunc 是一种函数类型，其格式为func(string)
type ListenFunc func(string)

func (lf ListenFunc) listen(inifile string) {
	lf(inifile)
}

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
