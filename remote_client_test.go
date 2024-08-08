/*
@Time : 2024/8/8 15:31
@Author : muyiming
@File : remote_client_test
@Software: GoLand
*/

package node_manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteClient(t *testing.T) {
	ssc := SSHConfig{
		IP:        "172.16.77.131",
		Port:      10086,
		SshUser:   "ttuser",
		SshPassWd: "123.com",
		SshKey:    "",
	}
	cli, err := NewRemoteClient(ssc)
	assert.Nil(t, err)
	sshCmd := fmt.Sprintf("touch testmim")
	output, err := cli.RunCommand(sshCmd)
	assert.Nil(t, err)
	fmt.Printf("%#v", output)
}
