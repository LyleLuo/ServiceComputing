package iniconfig

import (
	"fmt"
	"os"
)

func Perror(err error) {
	if err != nil {
		if _, err2 := fmt.Fprintln(os.Stderr, err); err2 != nil {
			panic(err2)
		}
		os.Exit(1)
	}
}
