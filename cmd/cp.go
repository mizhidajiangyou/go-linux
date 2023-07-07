package cmd

import (
	"io"
	"os"
	"path/filepath"
)

func CP(src, dst string) error {
	// 展开源文件或目录的协同变量
	src = os.ExpandEnv(src)
	dst = os.ExpandEnv(dst)

	// 获取源文件或目录的信息
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 如果源文件或目录是一个目录，则递归复制其子文件或子目录
	if info.IsDir() {
		return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 根据源文件或目录的路径和目标目录或文件的路径生成目标路径
			dstPath := filepath.Join(dst, path[len(src):])

			// 如果当前遍历到的是一个目录，则创建目录
			if info.IsDir() {
				return os.MkdirAll(dstPath, os.ModePerm)
			}

			// 如果当前遍历到的是一个文件，则复制文件
			return copyMyFile(path, dstPath)
		})
	}

	// 如果源文件或目录是一个文件，则直接复制文件
	return copyMyFile(src, dst)
}

func copyMyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer closeFile(srcFile)

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer closeFile(dstFile)

	// 将源文件内容复制到目标文件中
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 获取源文件的权限并设置目标文件的权限
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		// 处理错误
	}
}
