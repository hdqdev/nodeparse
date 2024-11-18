// pkg/utils/logger.go
package utils

import (
	"log"
	"path/filepath"
	"runtime"
)

func init() {
	// 设置日志格式，显示文件名和行号
	log.SetFlags(log.Ldate | log.Ltime)
}

func LogDebug(format string, v ...interface{}) {
	// 获取调用者的文件和行号
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	// 获取相对路径
	if rel, err := filepath.Rel(".", file); err == nil {
		file = rel
	}

	log.Printf("[DEBUG] %s:%d "+format, append([]interface{}{file, line}, v...)...)
}

func LogError(format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	if rel, err := filepath.Rel(".", file); err == nil {
		file = rel
	}

	log.Printf("[ERROR] %s:%d "+format, append([]interface{}{file, line}, v...)...)
}
