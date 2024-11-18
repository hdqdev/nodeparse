// internal/parser/ss.go
package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/deqdev/nodeparse/pkg/model"
	"github.com/deqdev/nodeparse/pkg/utils"
)

type SSParser struct{}

func (p *SSParser) CanParse(link string) bool {
	return strings.HasPrefix(link, "ss://")
}

func (p *SSParser) Parse(link string) (model.Node, error) {
	node := &model.SSNode{}

	// 1. 移除 "ss://" 前缀
	link = strings.TrimPrefix(link, "ss://")

	utils.LogDebug("开始解析 SS 链接: %s", link)

	// 2. 分离备注信息
	var mainPart string
	if idx := strings.Index(link, "#"); idx > -1 {
		mainPart = link[:idx]
		node.Name = link[idx+1:]
		if decoded, err := url.QueryUnescape(node.Name); err == nil {
			node.Name = decoded
		}
		utils.LogDebug("提取到备注: %s", node.Name)
	} else {
		mainPart = link
	}

	// 3. 尝试对整个主体部分进行 Base64 解码
	decodedBytes, err := base64.StdEncoding.DecodeString(mainPart)
	if err != nil {
		// 如果标准 Base64 解码失败，尝试 URL 安全的 Base64
		decodedBytes, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(mainPart)
		if err != nil {
			// 如果两种 Base64 都解码失败，假设是 SIP002 格式
			utils.LogDebug("整体 Base64 解码失败，尝试解析为 SIP002 格式")
			return p.parseSIP002(mainPart, node)
		}
	}

	content := string(decodedBytes)
	utils.LogDebug("Base64 解码后的内容: %s", content)

	// 4. 解析解码后的内容
	parts := strings.SplitN(content, "@", 2)
	if len(parts) != 2 {
		utils.LogError("链接格式错误：缺少 @ 分隔符")
		return nil, fmt.Errorf("missing @ separator in decoded content")
	}

	userInfo := parts[0]
	serverInfo := parts[1]

	// 5. 解析认证信息 (method:password)
	authParts := strings.SplitN(userInfo, ":", 2)
	if len(authParts) != 2 {
		utils.LogError("认证信息格式错误: %s", userInfo)
		return nil, fmt.Errorf("invalid method:password format: %s", userInfo)
	}
	node.Method = authParts[0]
	node.Password = authParts[1]

	// 6. 解析服务器信息 (server:port)
	serverParts := strings.SplitN(serverInfo, ":", 2)
	if len(serverParts) != 2 {
		utils.LogError("服务器信息格式错误: %s", serverInfo)
		return nil, fmt.Errorf("invalid server:port format: %s", serverInfo)
	}
	node.Server = serverParts[0]

	// 7. 解析端口
	port, err := strconv.Atoi(serverParts[1])
	if err != nil {
		utils.LogError("端口号格式错误: %s", serverParts[1])
		return nil, fmt.Errorf("invalid port number: %s", serverParts[1])
	}
	node.Port = port

	// 8. 验证必要字段
	if err := p.validateNode(node); err != nil {
		return nil, err
	}

	utils.LogDebug("SS 节点解析成功: %s@%s:%d", node.Method, node.Server, node.Port)

	return node, nil
}

func (p *SSParser) parseSIP002(mainPart string, node *model.SSNode) (model.Node, error) {
	// 处理 SIP002 格式的情况
	parts := strings.SplitN(mainPart, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SIP002 format: missing @")
	}

	// 解码用户信息部分
	userInfoBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info in SIP002 format: %v", err)
	}

	userInfo := string(userInfoBytes)
	authParts := strings.SplitN(userInfo, ":", 2)
	if len(authParts) != 2 {
		return nil, fmt.Errorf("invalid method:password format in SIP002")
	}

	node.Method = authParts[0]
	node.Password = authParts[1]

	// 处理服务器信息
	serverParts := strings.SplitN(parts[1], ":", 2)
	if len(serverParts) != 2 {
		return nil, fmt.Errorf("invalid server:port format in SIP002")
	}

	node.Server = serverParts[0]

	// 处理端口和可能的插件参数
	portAndParams := serverParts[1]
	if idx := strings.Index(portAndParams, "?"); idx != -1 {
		portStr := portAndParams[:idx]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port number in SIP002: %v", err)
		}
		node.Port = port

		// 处理插件参数
		params := portAndParams[idx+1:]
		if err := p.parsePluginParams(params, node); err != nil {
			return nil, err
		}
	} else {
		port, err := strconv.Atoi(portAndParams)
		if err != nil {
			return nil, fmt.Errorf("invalid port number in SIP002: %v", err)
		}
		node.Port = port
	}

	return node, p.validateNode(node)
}

func (p *SSParser) parsePluginParams(params string, node *model.SSNode) error {
	// 解析插件参数，如果有的话
	values, err := url.ParseQuery(params)
	if err != nil {
		return fmt.Errorf("failed to parse plugin params: %v", err)
	}

	if plugin := values.Get("plugin"); plugin != "" {
		// 处理插件信息
		utils.LogDebug("检测到插件: %s", plugin)
		// 这里可以添加插件相关的处理逻辑
	}

	return nil
}

func (p *SSParser) validateNode(node *model.SSNode) error {
	if node.Method == "" {
		utils.LogError("加密方法为空")
		return fmt.Errorf("empty encryption method")
	}
	if node.Password == "" {
		utils.LogError("密码为空")
		return fmt.Errorf("empty password")
	}
	if node.Server == "" {
		utils.LogError("服务器地址为空")
		return fmt.Errorf("empty server address")
	}
	if node.Port <= 0 || node.Port > 65535 {
		utils.LogError("端口号超出范围: %d", node.Port)
		return fmt.Errorf("invalid port range: %d", node.Port)
	}

	// 如果没有名称，使用服务器地址作为名称
	if node.Name == "" {
		node.Name = fmt.Sprintf("SS-%s:%d", node.Server, node.Port)
	}

	return nil
}
