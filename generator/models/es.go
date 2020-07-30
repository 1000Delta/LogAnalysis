package models

import (
	"fmt"
	"strings"
)

// ESCluster 定义 ES 集群信息
type ESCluster struct {
	Name  string
	Nodes []*ESNode
}

// NewESCluster 创建集群对象
func NewESCluster() *ESCluster {
	return &ESCluster{
		Name:  "",
		Nodes: []*ESNode{},
	}
}

// Hosts 返回集群中节点列表的字符串，参数为忽略节点的map
// ignore: 忽略节点列表
// 	example return: node1,node2,node3,node5
func (n *ESCluster) Hosts(ignore string) string {
	seeds := ""
	igNodes := strings.Split(ignore, ",")
	for _, node := range n.Nodes {
		exist := false
		for _, ig := range igNodes {
			if ig == node.Name {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		seeds += fmt.Sprintf("%s,", node.Name)
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
