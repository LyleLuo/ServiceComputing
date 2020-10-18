package main

import (
	"github.com/LyleLuo/ServiceComputing/pconfig/iniconfig"
)

func main() {
	c, err := iniconfig.InitConfig("test.ini")
	if err != nil {
		iniconfig.Perror(err)
	}
	iniconfig.Watch("test.ini", c)
}
