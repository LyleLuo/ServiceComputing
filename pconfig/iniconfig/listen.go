package iniconfig

import (
	"reflect"
	"time"
)

type Listener interface {
	listen(inifile string)
}

type ListenFunc func(string)

func (lf ListenFunc) listen(inifile string) {
	lf(inifile)
}

func (c *Config) ListenFunc(inifile string) {
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
