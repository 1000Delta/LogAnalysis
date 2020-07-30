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
	esc  = NewESCluster()
	once = &sync.Once{}
)

func singleton() {
	once.Do(func() {
		esc.Name = "test-cluster"
		for i := 1; i < 4; i++ {
			esc.Nodes = append(esc.Nodes, &ESNode{
				Name:     fmt.Sprintf("test-node-%d", i),
				IsMaster: true,
				HeapSize: 512,
			})
		}
	})
}

func TestSeedHosts(t *testing.T) {
	singleton()
	fail := func(tgt, now string) {
		t.Fatalf("输出匹配失败，目标输出: \"%v\", 实际输出: \"%v\"", tgt, now)
	}
	var hosts string
	// 不忽略节点
	hosts = esc.Hosts("")
	if hosts != InitialMasterNodes {
		fail(hosts, InitialMasterNodes)
	}
	// 忽略一个节点
	for i := 0; i < 3; i++ {
		hosts = esc.Hosts(esc.Nodes[i].Name)
		if hosts != seedHosts[i] {
			fail(seedHosts[i], hosts)
		}
	}
	// 忽略多个节点
	for i := 0; i < 3; i++ {
		n1, n2 := i, (i+1)%3
		hosts = esc.Hosts(fmt.Sprintf("%s,%s", esc.Nodes[n1].Name, esc.Nodes[n2].Name))
		if hosts != multiHosts[i] {
			fail(multiHosts[i], hosts)
		}
	}
}

func TestMasterNodes(t *testing.T) {
	singleton()
	nodes := esc.MasterNodes()
	if nodes != InitialMasterNodes {
		t.Fatalf("目标值: \"%s\", 当前值: \"%s\"", nodes, InitialMasterNodes)
	}
}
