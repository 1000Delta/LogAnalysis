package models

import "fmt"

// ESCluster 定义 ES 集群信息
type ESCluster struct {
	Name  string
	Nodes []*ESNode
}

// NewESCluster
func NewESCluster() *ESCluster {
	return &ESCluster{
		Name:  "",
		Nodes: []*ESNode{},
	}
}

// SeedHosts 返回集群中节点需要发现的节点
func (n *ESCluster) SeedHosts(node *ESNode) string {
	seeds := ""
	for _, n := range n.Nodes {
		if n == node {
			continue
		}
		seeds += fmt.Sprintf("%s,", n.Name)
	}
	return seeds[:len(seeds)-1]
}

// ESNode 描述 ES 节点信息
type ESNode struct {
	Name     string
	IsMaster bool
	// HeapSize 描述集群中节点 Java 堆大小，MiB
	HeapSize int
}
