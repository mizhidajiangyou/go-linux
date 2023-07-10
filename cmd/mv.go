package cmd

import (
	"os"
	"path/filepath"
)

func Mv(src, dst string) error {
	// 展开源文件或目录的环境变量
	src = os.ExpandEnv(src)

	// 解析源文件或目录的路径
	srcPath := filepath.Clean(src)

	// 解析目标目录或文件的路径
	dstPath := filepath.Clean(dst)

	// 获取源文件或目录的信息
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	// 获取目标路径的信息
	dstInfo, err := os.Stat(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 如果目标路径存在且是一个目录，则将源文件或目录移动到该目录下
	if dstInfo != nil && dstInfo.IsDir() {
		dstPath = filepath.Join(dstPath, filepath.Base(srcPath))
	}

	// 移动源文件或目录到目标路径
	err = os.Rename(srcPath, dstPath)
	if err != nil {
		return err
	}

	// 如果源文件或目录是一个目录，则递归移动其子文件或子目录
	if info.IsDir() {
		return filepath.Walk(dstPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 根据源文件或目录的路径和目标目录或文件的路径生成目标路径
			dstPath := filepath.Join(dstPath, path[len(srcPath):])

			// 如果当前遍历到的是一个目录，则创建目录
			if info.IsDir() {
				err = os.MkdirAll(dstPath, info.Mode())
				if err != nil {
					return err
				}
				return nil
			}

			// 如果当前遍历到的是一个文件，则移动文件
			return Mv(path, dstPath)
		})
	}

	return nil
}
