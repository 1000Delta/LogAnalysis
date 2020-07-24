package main

import "testing"

func TestReadConfig(t *testing.T) {
	// test json
	if _, err := readConfig("./test/test.json"); err != nil {
		t.Errorf("读取 json 错误：%v", err.Error())
	}
	// test yaml
	if _, err := readConfig("./test/test.yml"); err != nil {
		t.Errorf("读取 yaml 错误：%v", err.Error())
	}
	// test other
	if _, err := readConfig("./test/test.other"); err == nil {
		t.Errorf("读取不支持类型未报错")
	}
}

func TestGetPipelines(t *testing.T) {
	if _, err := GetPipelines(); err != nil {
		t.Errorf("加载 pipeline 错误：%v", err.Error())
	}
}
