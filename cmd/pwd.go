package cmd

import "os"

func Pwd() (string, error) {
	// 获取当前工作目录
	dir, err := os.Getwd()

	// 如果出现错误，返回错误信息
	if err != nil {
		return "", err
	}

	return dir, nil
}
