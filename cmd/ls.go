package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
)

func LsDir(path string) ([]string, error) {
	// 将路径中的环境变量替换为其对应的值
	path = os.ExpandEnv(path)

	// 获取指定路径的文件信息
	info, err := os.Stat(path)

	// 如果出现错误，返回错误信息
	if err != nil {
		return nil, err
	}

	// 如果指定的是目录，列出目录下的文件和子目录
	if info.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}

		// 构造文件名切片
		names := make([]string, len(files))
		for i, f := range files {
			names[i] = f.Name()
		}

		return names, nil
	} else {
		// 如果指定的是文件，返回错误信息
		return nil, os.ErrNotExist
	}
}

func LsFile(path string) error {
	// 将路径中的环境变量替换为其对应的值
	path = os.ExpandEnv(path)

	// 获取指定路径的文件信息
	info, err := os.Stat(path)

	// 如果出现错误，返回错误信息
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("something went wrong: %s", "missing file")
	} else {
		return nil
	}
}
