package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Grep grep
func Grep(filename, searchText string) (s []string, err error) {
	// 判断文件字符串中是否包含环境变量
	if strings.Contains(filename, "$") {
		// 将文件中的环境变量替换为其对应的值
		filename = os.ExpandEnv(filename)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return s, fmt.Errorf("%s does not exist", filename)
	}

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return s, err
	}
	defer file.Close()

	// 创建一个带缓冲区的文件读取器
	reader := bufio.NewReader(file)

	// 用于记录行号
	lineNumber := 0

	// 循环读取文件中的每一行
	for {
		line, err := reader.ReadString('\n')

		// 如果已经读取到文件末尾，则退出循环
		if err != nil && err.Error() == "EOF" {
			break
		}

		// 如果读取出现错误，则返回错误信息
		if err != nil {
			return s, err
		}

		// 增加行号
		lineNumber++

		// 如果行中包含搜索字符串，则写入标准输出流
		if strings.Contains(line, searchText) {
			s = append(s, fmt.Sprintf("%s:%d:%s", filename, lineNumber, line))
		}
	}

	return s, nil
}

// VGrep grep -v
func VGrep(filename, searchText string) (s []string, err error) {
	// 判断文件字符串中是否包含环境变量
	if strings.Contains(filename, "$") {
		// 将文件中的环境变量替换为其对应的值
		filename = os.ExpandEnv(filename)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return s, fmt.Errorf("%s does not exist", filename)
	}

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return s, err
	}
	defer file.Close()

	// 创建一个带缓冲区的文件读取器
	reader := bufio.NewReader(file)

	// 用于记录行号
	lineNumber := 0

	// 循环读取文件中的每一行
	for {
		line, err := reader.ReadString('\n')

		// 如果已经读取到文件末尾，则退出循环
		if err != nil && err.Error() == "EOF" {
			break
		}

		// 如果读取出现错误，则返回错误信息
		if err != nil {
			return s, err
		}

		// 增加行号
		lineNumber++

		// 如果行中不包含搜索字符串，则写入标准输出流
		if !strings.Contains(line, searchText) {
			s = append(s, fmt.Sprintf("%s:%d:%s", filename, lineNumber, line))
		}
	}

	return s, nil
}
