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

// ErrConfigNotExist 错误：配置列表还没有被初始化
var ErrConfigNotExist = errors.New("The Config struct haven't been constructed")

func init() {
	if runtime.GOOS == "windows" {
		commentChar = ';'
	} else {
		commentChar = '#'
	}
}

// Config 储存配置文件的信息
type Config struct {
	filepath string
	Conflist map[string]map[string]string
}

// InitConfig 读取配置文件的信息
func InitConfig(filepath string) (*Config, error) {
	c := new(Config)
	c.filepath = filepath
	err := c.readList()
	return c, err
}

// GetValue 通过section和name获得值
func (c *Config) GetValue(section, name string) (string, error) {
	if _, ok := c.Conflist[section][name]; ok {
		return c.Conflist[section][name], nil
	}
	return "", ErrKeyNotExist
}

// readList 读取配置文件
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
	for {
		l, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
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
			i := strings.IndexAny(line, "=")
			if i == -1 {
				continue
			}
			value := strings.TrimSpace(line[i+1 : len(line)])
			sectionMap[strings.TrimSpace(line[0:i])] = value
		}
	}

	c.Conflist[section] = sectionMap
	return nil
}

// GetAllSetion 返回所有的section和其中的键值对
func (c *Config) GetAllSetion() (map[string]map[string]string, error) {
	if c.Conflist == nil {
		return nil, ErrConfigNotExist
	}
	return c.Conflist, nil
}
