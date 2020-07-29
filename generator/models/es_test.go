package models

import (
	"fmt"
	"testing"
)

var (
	seedHosts = []string{
		"test-node-2,test-node-3",
		"test-node-1,test-node-3",
		"test-node-1,test-node-2",
	}
)

func TestSeedHosts(t *testing.T) {
	esc := NewESCluster()
	esc.Name = "test-cluster"
	for i := 1; i < 4; i++ {
		esc.Nodes = append(esc.Nodes, &ESNode{Name: fmt.Sprintf("test-node-%d", i)})
	}
	for i := 0; i < 3; i++ {
		if (esc.SeedHosts(esc.Nodes[i]) != seedHosts[i]) {
			t.Fatalf("输出匹配失败，格式输出: \"%v\", 实际输出: \"%v\"", seedHosts[i], esc.SeedHosts(esc.Nodes[i]))
		}
	}
}
