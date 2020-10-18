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
