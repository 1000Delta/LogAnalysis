package models

import (
	"fmt"
	"testing"
)

func TestSeedHosts(t *testing.T) {
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

		fail = func(tgt, now string) {
			t.Fatalf("输出匹配失败，目标输出: \"%v\", 实际输出: \"%v\"", tgt, now)
		}
	)
	esc := NewESCluster()
	esc.Name = "test-cluster"
	for i := 1; i < 4; i++ {
		esc.Nodes = append(esc.Nodes, &ESNode{Name: fmt.Sprintf("test-node-%d", i)})
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
