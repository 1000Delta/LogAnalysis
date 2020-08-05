package main

import (
	"fmt"
	"generator/models"
	"os"
	"sync"
	"testing"
	"text/template"
)

var (
	esc  = models.NewESCluster("test-cluster")
	once = &sync.Once{}
)

func singleton() {
	once.Do(func() {
		esc.Nodes = append(esc.Nodes, models.NewESNode(
			fmt.Sprintf("test-node-%d", 1),
			map[int]int{9200: 9200},
			512,
			true,
			fmt.Sprintf("data%02d", 1),
		))
		for i := 2; i < 4; i++ {
			esc.Nodes = append(esc.Nodes, models.NewESNode(
				fmt.Sprintf("test-node-%d", i),
				nil,
				512,
				true,
				fmt.Sprintf("data%02d", i),
			))
		}
	})
}

func TestTemplateFormat(t *testing.T) {
	singleton()
	tmplES, err := template.ParseFiles("./templates/docker-compose.yml.tmpl", "./templates/docker-compose.yml.d/es.tmpl")
	if err != nil {
		t.Fatalf("解析 es 模板错误：%s", err.Error())
	}
	outES, err := os.Create("./build/docker-compose.yml")
	if err != nil {
		t.Fatalf("创建 es 文件失败：%s", err.Error())
	}
	defer outES.Close()
	if err := tmplES.Execute(outES, esc); err != nil {
		t.Fatalf("模板执行失败：%s", err.Error())
	}
}
