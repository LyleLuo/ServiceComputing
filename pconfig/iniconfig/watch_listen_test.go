package iniconfig

import "testing"

// 运行此测试的过程中把http_port = 9999改为http_port = 8888
func TestWatch(t *testing.T) {
	var f ListenFunc = Listen
	m, err := Watch("../test.ini", f)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if m[""]["app_mode"] != "development" {
		t.Errorf("Not equal: expected: development, actual: %s", m[""]["app_mode"])
	}
	if m["paths"]["data"] != "/home/git/grafana" {
		t.Errorf("Not equal: expected: /home/git/grafana, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["protocol"] != "http" {
		t.Errorf("Not equal: expected: http, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["http_port"] != "8888" {
		t.Errorf("Not equal: expected: 8888, actual: %s", m[""]["app_mode"])
	}
	if m["server"]["enforce_domain"] != "true" {
		t.Errorf("Not equal: expected: true, actual: %s", m[""]["app_mode"])
	}
}
