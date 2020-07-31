package models

import (
	"fmt"
	"sync"
	"testing"
)

var (
	InitialMasterNodes = "test-node-1,test-node-2,test-node-3"
	seedHosts          = []string{
		"test-node-2,test-node-3",
		"test-node-1,test-node-3",
		"test-node-1,test-node-2",
	}
	multiHosts = []string{
		"test-node-3",
		"test-node-1",
		"test-node-2",
	}
)

var (
	esc  = NewESCluster("test-cluster")
	once = &sync.Once{}
)

func singleton() {
	once.Do(func() {
		esc.Nodes = append(esc.Nodes, NewESNode(
			fmt.Sprintf("test-node-%d", 1),
			map[int]int{9200: 9200},
			512,
			true,
			fmt.Sprintf("data%02d", 1),
		))
		for i := 2; i < 4; i++ {
			esc.Nodes = append(esc.Nodes, NewESNode(
				fmt.Sprintf("test-node-%d", i),
				nil,
				512,
				true,
				fmt.Sprintf("data%02d", i),
			))
		}
	})
}

func TestGetHosts(t *testing.T) {
	singleton()
	var hosts string
	// 不忽略节点
	hosts = esc.GetHosts("")
	if hosts != InitialMasterNodes {
		t.Fatalf("输出匹配失败，目标输出: \"%v\", 实际输出: \"%v\"", hosts, InitialMasterNodes)
	}
	// 忽略一个节点
	for i := 0; i < 3; i++ {
		hosts = esc.GetHosts(esc.Nodes[i].Name)
		if hosts != seedHosts[i] {
			t.Fatalf("输出匹配失败，目标输出: \"%v\", 实际输出: \"%v\"", seedHosts[i], hosts)
		}
	}
	// 忽略多个节点
	for i := 0; i < 3; i++ {
		n1, n2 := i, (i+1)%3
		hosts = esc.GetHosts(fmt.Sprintf("%s,%s", esc.Nodes[n1].Name, esc.Nodes[n2].Name))
		if hosts != multiHosts[i] {
			t.Fatalf("输出匹配失败，目标输出: \"%v\", 实际输出: \"%v\"", multiHosts[i], hosts)
		}
	}
}

func TestMasterNodes(t *testing.T) {
	singleton()
	nodes := esc.GetMasterNodes()
	if nodes != InitialMasterNodes {
		t.Fatalf("目标值: \"%s\", 当前值: \"%s\"", nodes, InitialMasterNodes)
	}
}
