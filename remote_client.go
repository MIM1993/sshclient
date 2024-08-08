package node_manager

/*
	使用golang代码进行ssh连接服务器，并执行命令，长传文件
*/

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"sync"
	"time"
)

type RemoteClient struct {
	Client *ssh.Client
	Config *ssh.ClientConfig
	mu     sync.Mutex
}

type SSHConfig struct {
	IP        string `json:"ip"`         //ip
	Port      int    `json:"port"`       //端口
	SshUser   string `json:"ssh_user"`   //ssh用户名
	SshPassWd string `json:"ssh_passwd"` //ssh密码
	SshKey    string `json:"ssh_key"`    //ssh免密文件
}

func NewRemoteClient(sshConfig SSHConfig) (client *RemoteClient, err error) {
	hostIPConfig := fmt.Sprintf("%s:%d", sshConfig.IP, sshConfig.Port)
	client = &RemoteClient{
		Config: &ssh.ClientConfig{
			User: sshConfig.SshUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(sshConfig.SshPassWd),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产中应更安全地处理HostKey
			Timeout:         time.Second * 2,
		},
	}
	if err = client.connect(hostIPConfig); err != nil {
		return client, err
	}
	return client, err
}

// connect 是实际建立SSH连接的内部方法
func (ops *RemoteClient) connect(host string) error {
	var (
		client *ssh.Client
		err    error
	)
	for i := 0; i < 1; i++ {
		client, err = ssh.Dial("tcp", host, ops.Config)
		if err == nil {
			ops.Client = client
			return nil
		}
	}
	if err != nil {
		return err
	}
	return err
}

// RunCommand 在远程Linux系统上执行命令
func (ops *RemoteClient) RunCommand(command string, args ...string) (string, error) {
	ops.mu.Lock()
	defer ops.mu.Unlock()
	session, err := ops.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout
	cmd := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return stdout.String(), nil
}

// TransferFile 通过SCP将文件从本地传输到远程系统
func (ops *RemoteClient) TransferFile(localFile io.Reader, remotePath string, isExecutable bool) error {
	ops.mu.Lock()
	defer ops.mu.Unlock()
	sftpClient, err := sftp.NewClient(ops.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, localFile)
	if isExecutable {
		sftpClient.Chmod(remotePath, 0755)
	}

	return err
}

func (ops *RemoteClient) ClientClose() error {
	return ops.Client.Close()
}
