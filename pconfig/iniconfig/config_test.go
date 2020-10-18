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
