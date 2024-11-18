package model

type SSNode struct {
	BaseNode
	Password string `json:"password"`
	Method   string `json:"method"`
	Plugin   string `json:"plugin,omitempty"`
}

// 实现 Parse 方法
func (n *SSNode) Parse(link string) error {
	// 实现解析逻辑
	return nil
}

func (n *SSNode) GetType() string {
	return "ss"
}

func (n *SSNode) ToClashConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":     n.Name,
		"type":     "ss",
		"server":   n.Server,
		"port":     n.Port,
		"password": n.Password,
		"cipher":   n.Method,
	}
}
