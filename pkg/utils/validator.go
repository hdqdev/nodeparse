// pkg/utils/validator.go
package utils

// 添加一些基本的验证函数
func ValidatePort(port int) bool {
	return port > 0 && port <= 65535
}

func ValidateServerAddress(address string) bool {
	return len(address) > 0
}

func ValidateShadowsocksMethod(method string) bool {
	validMethods := map[string]bool{
		"aes-128-gcm":             true,
		"aes-192-gcm":             true,
		"aes-256-gcm":             true,
		"chacha20-poly1305":       true,
		"xchacha20-ietf-poly1305": true,
	}
	return validMethods[method]
}

func ValidatePassword(password string) bool {
	return len(password) > 0
}
