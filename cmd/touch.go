package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Touch(files ...string) error {
	for _, file := range files {
		// 判断文件字符串中是否包含环境变量
		if strings.Contains(file, "$") {
			// 将文件中的环境变量替换为其对应的值
			file = os.ExpandEnv(file)
		}

		// 获取目录路径
		dir := filepath.Dir(file)

		// 如果目录不存在，则创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// 更新文件时间戳
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// 如果文件不存在，则创建一个空文件
			if _, err := os.Create(file); err != nil {
				return err
			}
		} else {
			// 如果文件存在，则更新时间戳
			if err := os.Chtimes(file, time.Now(), time.Now()); err != nil {
				return err
			}
		}
	}
	return nil
}
