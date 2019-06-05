package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/kevinburke/ssh_config"
	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/common/log"
	"golang.org/x/crypto/ssh"
)

func publicKey(path string) ssh.AuthMethod {
	fullpath, err := homedir.Expand(path)
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadFile(fullpath)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		panic(err)
	}
}

func createConnection(alias string) *ssh.Client {
	config := &ssh.ClientConfig{
		User: ssh_config.Get(alias, "User"),
		Auth: []ssh.AuthMethod{
			publicKey(ssh_config.Get(alias, "IdentityFile")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", ssh_config.Get(alias, "Hostname")+":"+ssh_config.Get(alias, "Port"), config)
	if err != nil {
		log.Errorln(err)
	}
	return conn
}

func main() {
	conn := createConnection("example")
	runCommand("whoami", conn)
	defer conn.Close()
}
