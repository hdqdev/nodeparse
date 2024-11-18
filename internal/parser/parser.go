package parser

import (
	"errors"

	"github.com/deqdev/nodeparse/pkg/model"
)

type Parser interface {
	Parse(link string) (model.Node, error)
	CanParse(link string) bool
}

// 常见错误定义
var (
	ErrInvalidLink         = errors.New("invalid link format")
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
)

// NodeParser 用于管理所有协议的解析器
type NodeParser struct {
	parsers []Parser
}

// NewNodeParser 创建一个新的节点解析器
func NewNodeParser() *NodeParser {
	return &NodeParser{
		parsers: []Parser{
			&SSParser{},
			// &VMessParser{},
			// 在这里添加其他协议的解析器
		},
	}
}

// ParseLink 解析单个链接
func (np *NodeParser) ParseLink(link string) (model.Node, error) {
	for _, parser := range np.parsers {
		if parser.CanParse(link) {
			return parser.Parse(link)
		}
	}
	return nil, ErrUnsupportedProtocol
}

// ParseLinks 批量解析多个链接
func (np *NodeParser) ParseLinks(links []string) ([]model.Node, error) {
	nodes := make([]model.Node, 0, len(links))
	errors := make([]error, 0)

	for _, link := range links {
		node, err := np.ParseLink(link)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		nodes = append(nodes, node)
	}

	if len(nodes) == 0 && len(errors) > 0 {
		return nil, errors[0] // 返回第一个错误
	}

	return nodes, nil
}
