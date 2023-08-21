package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func Mkdir(dir string) error {
	// 判断目录字符串中是否包含环境变量
	if strings.Contains(dir, "$") {
		// 将目录中的环境变量替换为其对应的值
		dir = os.ExpandEnv(dir)
	}
	// 创建目录
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		// 如果创建目录失败，则尝试单个目录创建
		parts := strings.Split(dir, "/")
		for i := range parts {
			path := filepath.Join(parts[:i+1]...)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if err := os.Mkdir(path, os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
