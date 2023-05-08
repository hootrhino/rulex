package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

func Test_ssh_tty(t *testing.T) {
	sshHost := "192.168.0.105"
	sshUser := "ubuntu"
	sshPassword := "ubuntupi"
	sshType := "password" //password 或者 key
	sshKeyPath := ""      //ssh id_rsa.id 路径"
	sshPort := 22

	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		glogger.GLogger.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		glogger.GLogger.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令
	combo, err := session.CombinedOutput("whoami; cd /; ls -al;echo https://github.com/dejavuzhou/felix")
	if err != nil {
		glogger.GLogger.Fatal("远程执行cmd 失败", err)
	}
	glogger.GLogger.Println("命令输出:", string(combo))

}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		glogger.GLogger.Fatal("find key's home dir failed", err)
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		glogger.GLogger.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		glogger.GLogger.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
