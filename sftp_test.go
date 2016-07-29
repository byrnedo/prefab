package prefab

import (
	"testing"
	"time"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func getSftpClient(addr string, t *testing.T) *sftp.Client {

	sshConfig := &ssh.ClientConfig{
		User: SftpUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(SftpPassword),
		},
	}

	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		t.Fatal(err)
	}

	sftpClient, err := sftp.NewClient(sshClient, func(*sftp.Client)error{return nil})
	if err != nil {
		t.Fatal(err)
	}
	return sftpClient
}

func TestSetupSftpContainer(t *testing.T) {
	id, url := StartSftpContainer()
	if url == "" {
		t.Fatal("url empty")
	}
	if id == "" {
		t.Fatal("id empty")
	}

	if err := WaitForSftp(url, 30 * time.Second); err != nil {
		t.Fatal(err.Error())
	}

	client := getSftpClient(url, t)

	_, err := client.ReadDir("./data")
	if err != nil {
		t.Fatal(err)
	}

	Remove(id)
}
