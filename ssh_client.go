package node_manager

import (
	"io"
)

type NodeClient interface {
	RunCommand(command string, args ...string) (string, error)
	TransferFile(localFile io.Reader, remotePath string, isExecutable bool) error
	ClientClose() error
}
