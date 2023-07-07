package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Scp(src, dest string) error {
	srcPath, srcHost, err := parsePath(src)
	if err != nil {
		return err
	}

	destPath, destHost, err := parsePath(dest)
	if err != nil {
		return err
	}

	srcClient, err := getSftpClient(srcHost)
	if err != nil {
		return err
	}
	defer srcClient.Close()

	destClient, err := getSftpClient(destHost)
	if err != nil {
		return err
	}
	defer destClient.Close()

	return copy(srcClient, destClient, srcPath, destPath)
}

func parsePath(input string) (string, string, error) {
	parts := strings.Split(input, ":")
	if len(parts) == 1 {
		return os.ExpandEnv(input), "", nil
	} else if len(parts) == 2 {
		return parts[1], parts[0], nil
	}
	return "", "", fmt.Errorf("invalid input format")
}

func getSftpClient(host string) (*sftp.Client, error) {
	if host == "" {
		return sftp.NewClient(nil)
	}

	sshConfig := &ssh.ClientConfig{
		User:            os.Getenv("USER"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("PASSWORD")),
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), sshConfig)
	if err != nil {
		return nil, err
	}

	return sftp.NewClient(conn)
}

func copy(srcClient, destClient *sftp.Client, srcPath, destPath string) error {
	srcFile, err := srcClient.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if srcFileInfo.IsDir() {
		return copyDir(srcClient, destClient, srcPath, destPath)
	}

	return copyFile(srcClient, destClient, srcPath, destPath)
}

func copyFile(srcClient, destClient *sftp.Client, srcPath, destPath string) error {
	srcFile, err := srcClient.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := destClient.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = srcFile.WriteTo(destFile)
	return err
}

func copyDir(srcClient, destClient *sftp.Client, srcPath, destPath string) error {
	srcFiles, err := srcClient.ReadDir(srcPath)
	if err != nil {
		return err
	}

	err = destClient.Mkdir(destPath)
	if err != nil && !os.IsExist(err) {
		return err
	}

	for _, srcFile := range srcFiles {
		srcFilePath := filepath.Join(srcPath, srcFile.Name())
		destFilePath := filepath.Join(destPath, srcFile.Name())

		if srcFile.IsDir() {
			err = copyDir(srcClient, destClient, srcFilePath, destFilePath)
		} else {
			err = copyFile(srcClient, destClient, srcFilePath, destFilePath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
