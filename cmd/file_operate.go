package cmd

import (
	"bufio"
	"os"
	"strings"
)

// InsertData 在文件的指定行或首尾插入数据。
// 如果 lineNum 为 0，则在文件首部插入数据。
// 如果 lineNum 为 -1，则在文件尾部插入数据。
func InsertData(filename, data string, lineNum int) error {
	// 替换文件路径中的环境变量
	if strings.Contains(filename, "$") {
		filename = os.ExpandEnv(filename)
	}

	// 打开文件
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个 Scanner，逐行读取文件
	scanner := bufio.NewScanner(file)

	// 创建一个缓冲区，用于存储新的文件内容
	var buffer []byte

	// 记录当前行号
	currentLineNumber := 0

	// 计算要插入数据的位置
	insertIndex := lineNum - 1
	if insertIndex < 0 {
		insertIndex = 0
	}

	// 循环遍历文件的每一行
	for scanner.Scan() {
		// 增加当前行号
		currentLineNumber++

		// 如果当前行号等于要插入的行号，则将新数据添加到缓冲区中
		if currentLineNumber == lineNum {
			buffer = append(buffer, []byte(data)...)
			buffer = append(buffer, '\n')
		}

		// 将当前行添加到缓冲区中
		buffer = append(buffer, scanner.Bytes()...)
		buffer = append(buffer, '\n')

		// 如果当前行号等于要插入的行号，并且要插入数据的位置不是文件末尾，则将要插入的数据插入到缓冲区中
		if currentLineNumber == lineNum && insertIndex != len(buffer) {
			buffer = append(buffer[:insertIndex], append([]byte(data+"\n"), buffer[insertIndex:]...)...)
		}
	}

	// 如果要插入的行号大于文件中的行数，则在文件末尾插入新数据
	if lineNum > currentLineNumber {
		buffer = append(buffer, []byte(data)...)
		buffer = append(buffer, '\n')
	}

	// 如果要插入的行号为 0，则在文件首部插入新数据
	if lineNum == 0 {
		buffer = append([]byte(data), '\n')
		buffer = append(buffer, buffer...)
	}

	// 如果要插入的行号为 -1，则在文件尾部插入新数据
	if lineNum == -1 {
		buffer = append(buffer, []byte(data)...)
		buffer = append(buffer, '\n')
	}

	// 截断文件为零长度
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// 将文件指针移动到文件开头
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// 将新的文件内容写入文件
	_, err = file.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}
