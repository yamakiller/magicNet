package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDirector : 获取当前目录
func GetCurrentDirector() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
