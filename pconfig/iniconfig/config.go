// Package iniconfig 用于读取ini配置文件的包
package iniconfig

import (
	"bufio"
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
)

var commentChar rune

// ErrKeyNotExist 错误：不存在该键值对
var ErrKeyNotExist = errors.New("The section don't include this key-value pair")

func init() {
	if runtime.GOOS == "windows" {
		commentChar = ';'
	} else {
		commentChar = '#'
	}
}

// Config 储存配置文件的信息，包括文件名和文件里的键值对
type Config struct {
	filepath string
	Conflist map[string]map[string]string
}

// InitConfig 读取配置文件的信息，参数filepath是要读取的文件路径，返回值*Config是存储文件信息的结构体，返回值err是错误类型
func InitConfig(filepath string) (*Config, error) {
	c := new(Config)
	c.filepath = filepath
	err := c.readList()
	return c, err
}

// GetValue 通过参数section和name获得值，返回string是value，error是错误类型
func (c *Config) GetValue(section, name string) (string, error) {
	if _, ok := c.Conflist[section][name]; ok {
		return c.Conflist[section][name], nil
	}
	return "", ErrKeyNotExist
}

// readList 读取配置文件的方法，返回值是错误类型
func (c *Config) readList() error {
	file, err := os.Open(c.filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	c.Conflist = make(map[string]map[string]string)
	var section string
	var sectionMap = make(map[string]string)
	isFirstSection := true
	buf := bufio.NewReader(file)
	iter := 1
	for {
		l, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line := strings.TrimSpace(l)
		switch {
		case len(line) == 0:
		case string(line[0]) == string(commentChar):
		case line[0] == '[' && line[len(line)-1] == ']':
			if !isFirstSection {
				c.Conflist[section] = sectionMap
			} else {
				isFirstSection = false
			}
			section = strings.TrimSpace(line[1 : len(line)-1])
			sectionMap = make(map[string]string)
		default:
			if iter == 1 {
				isFirstSection = false
			}
			i := strings.IndexAny(line, "=")
			if i == -1 {
				continue
			}
			value := strings.TrimSpace(line[i+1 : len(line)])
			sectionMap[strings.TrimSpace(line[0:i])] = value
		}
		if err == io.EOF {
			break
		}
	}

	c.Conflist[section] = sectionMap
	return nil
}
