package main

import (
	"fmt"

	"github.com/hauke-d/external-dns-sftp-webhook/transport"
	"golang.org/x/crypto/ssh"
)

func main() {
	transport, err := transport.NewSshTransport("localhost", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("root")},
	})
	if err != nil {
		panic(err)
	}

	transport.ReadFile("/etc/os-release")

	fmt.Println("Working")
}
