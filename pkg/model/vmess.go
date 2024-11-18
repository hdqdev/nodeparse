package model

type VmessNode struct {
	BaseNode
	UUID           string            `json:"uuid"`
	AlterID        int               `json:"alterId"`
	Security       string            `json:"security"`
	Network        string            `json:"network"`
	WSPath         string            `json:"ws-path,omitempty"`
	WSHeaders      map[string]string `json:"ws-headers,omitempty"`
	TLS            bool              `json:"tls"`
	SkipCertVerify bool              `json:"skip-cert-verify"`
}

func (n *VmessNode) GetType() string {
	return "vmess"
}

func (n *VmessNode) ToClashConfig() map[string]interface{} {
	// 实现 Vmess 的 Clash 配置转换
	return map[string]interface{}{
		"name":             n.Name,
		"type":             "vmess",
		"server":           n.Server,
		"port":             n.Port,
		"uuid":             n.UUID,
		"alterId":          n.AlterID,
		"cipher":           n.Security,
		"network":          n.Network,
		"ws-path":          n.WSPath,
		"ws-headers":       n.WSHeaders,
		"tls":              n.TLS,
		"skip-cert-verify": n.SkipCertVerify,
	}
}
