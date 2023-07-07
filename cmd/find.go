package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Find(name string) ([]string, error) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 递归查找文件
	var result []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，继续递归查找
		if info.IsDir() {
			return nil
		}

		// 如果文件名匹配，加入结果列表
		if strings.Contains(info.Name(), name) {
			result = append(result, path)
		}

		return nil
	})

	if len(result) == 0 {
		// 如果结果列表为空，返回错误信息
		if err == nil {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return result, err
}

func FindWithDir(name string, dir string) ([]string, error) {
	// 将目录中的环境变量替换为其对应的值
	dir = os.ExpandEnv(dir)

	// 递归查找文件
	var result []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，继续递归查找
		if info.IsDir() {
			return nil
		}

		// 如果文件名匹配，加入结果列表
		if strings.Contains(info.Name(), name) {
			result = append(result, path)
		}

		return nil
	})

	if len(result) == 0 {
		// 如果结果列表为空，返回错误信息
		if err == nil {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return result, err
}
