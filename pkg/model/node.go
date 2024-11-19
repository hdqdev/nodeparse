package model

// Node 节点接口
type Node interface {
	Parse(link string) error
	ToClashConfig() map[string]interface{}
	GetName() string
	GetType() string
	GetServer() string
	GetPort() int
}

// BaseNode 基础节点结构
type BaseNode struct {
	Name   string `json:"name"`
	Server string `json:"server"`
	Port   int    `json:"port"`
}

func (b *BaseNode) GetName() string {
	return b.Name
}

func (b *BaseNode) GetServer() string {
	return b.Server
}

func (b *BaseNode) GetPort() int {
	return b.Port
}
