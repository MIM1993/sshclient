package node_manager

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type LocalClient struct {
	mu sync.Mutex
}

func NewLocalClient() *LocalClient {
	return &LocalClient{}
}

// RunCommand 在远程Linux系统上执行命令
func (ops *LocalClient) RunCommand(command string, args ...string) (string, error) {
	ops.mu.Lock()
	defer ops.mu.Unlock()
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// TransferFile 通过SCP将文件从本地传输到远程系统
func (ops *LocalClient) TransferFile(localFile io.Reader, remotePath string, isExecutable bool) error {
	ops.mu.Lock()
	defer ops.mu.Unlock()

	dstFile, err := os.Create(remotePath)
	if err != nil {
		fmt.Printf("TransferFile dest: %s got error: %v", remotePath, err)
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, localFile)
	if isExecutable {
		err = os.Chmod(remotePath, 0755)
	}

	return err
}

func (ops *LocalClient) ClientClose() error {
	return nil
}
