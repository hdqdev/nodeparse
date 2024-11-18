package utils

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64Safe 解码 Base64，支持标准和 URL 安全的格式
func DecodeBase64Safe(s string) ([]byte, error) {
	// 处理填充
	if pad := len(s) % 4; pad > 0 {
		s += strings.Repeat("=", 4-pad)
	}

	// 尝试标准 Base64
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return decoded, nil
	}

	// 尝试 URL 安全的 Base64
	return base64.URLEncoding.DecodeString(s)
}
