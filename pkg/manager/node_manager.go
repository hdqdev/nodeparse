package manager

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hdqdev/nodeparse/internal/parser"
	"github.com/hdqdev/nodeparse/pkg/model"
	"github.com/hdqdev/nodeparse/pkg/utils"
)

type NodeManager struct {
	parsers []parser.Parser
	nodes   []model.Node
	options *Options
}

func NewNodeManager(opts ...Option) *NodeManager {
	nm := NodeManager{
		parsers: []parser.Parser{
			&parser.SSParser{},
			// &parser.VmessParser{},
		},
		options: &Options{},
	}
	for _, opt := range opts {
		opt(nm.options)
	}

	return &nm
}

func (nm *NodeManager) AddParser(parser parser.Parser) {
	nm.parsers = append(nm.parsers, parser)
}

// LoadFromURL 从 URL 加载节点配置
func (nm *NodeManager) LoadFromURL(url string) error {
	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发起 HTTP GET 请求
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析内容
	return nm.parseContent(string(body))
}

// LoadFromFile 从文件加载节点配置
func (nm *NodeManager) LoadFromFile(filename string) error {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 解析内容
	return nm.parseContent(string(content))
}

// parseContent 解析配置内容
func (nm *NodeManager) parseContent(content string) error {
	// 清理现有节点（可选，取决于你的需求）
	// nm.nodes = make([]model.Node, 0)

	// 按行分割内容
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		nodeParser := parser.NewNodeParser()
		node, err := nodeParser.ParseLink(line)
		if err != nil {
			utils.LogError("解析第 %d 行失败: %v", lineNum, err)
			// 可以选择继续解析，而不是直接返回错误
			continue
		}

		// 添加节点
		if err := nm.AddNode(node); err != nil {
			utils.LogError("添加节点失败: %v", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading content: %v", err)
	}

	utils.LogDebug("成功加载 %d 个节点", len(nm.nodes))
	return nil
}

// AddNode 添加节点到管理器
func (nm *NodeManager) AddNode(node model.Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	nm.nodes = append(nm.nodes, node)
	return nil
}

func (nm *NodeManager) GetNodes() []model.Node {
	return nm.nodes
}

func (nm *NodeManager) ExportToClash() []map[string]interface{} {
	var configs []map[string]interface{}
	for _, node := range nm.nodes {
		configs = append(configs, node.ToClashConfig())
	}
	return configs
}

func (nm *NodeManager) string() string {
	var sb strings.Builder
	for _, node := range nm.nodes {
		sb.WriteString(fmt.Sprintf("node: %v\n", node))
	}
	return sb.String()
}
