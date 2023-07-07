package cmd

import "os"

func Cd(dir string) error {
	// 将目录中的环境变量替换为其对应的值
	dir = os.ExpandEnv(dir)
	// 切换到指定目录
	err := os.Chdir(dir)

	return err
}
