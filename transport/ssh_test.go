package transport

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/ssh"
)

func connect(ctx context.Context, container testcontainers.Container, t *testing.T) Transport {
	host := getHost(ctx, container)
	transport, err := NewSshTransport(host, &ssh.ClientConfig{
		User:            "webhook",
		Auth:            []ssh.AuthMethod{ssh.Password("webhook")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		t.Fatal(err)
	}
	return transport
}

func startSshServer(ctx context.Context, t *testing.T) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Name:         "openssh-server",
		Image:        "lscr.io/linuxserver/openssh-server:latest",
		ExposedPorts: []string{"2222/tcp"},
		WaitingFor:   wait.ForLog("[ls.io-init] done."),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return container
}

func getHost(ctx context.Context, container testcontainers.Container) string {
	p, _ := nat.NewPort("tcp", "2222")
	port, _ := container.MappedPort(ctx, p)
	return fmt.Sprintf("localhost:%s", port.Port())
}

func TestRunCommand(t *testing.T) {
	ctx := context.Background()
	container := startSshServer(ctx, t)
	transport := connect(ctx, container, t)

	output, err := transport.Run("cat /etc/os-release")
	if err != nil {
		t.Fatal(err)
	}

	strings.Contains(string(output[:]), "Alpine Linux")
}

func TestReadFile(t *testing.T) {
	ctx := context.Background()
	container := startSshServer(ctx, t)
	transport := connect(ctx, container, t)

	contents, err := transport.ReadFile("/etc/os-release")
	if err != nil {
		t.Fatal(err)
	}

	strings.Contains(string(contents[:]), "Alpine Linux")
}

func TestWriteFile(t *testing.T) {
	ctx := context.Background()
	container := startSshServer(ctx, t)
	transport := connect(ctx, container, t)

	err := transport.WriteFile("new-os-release", []byte("SSH Linux"))
	if err != nil {
		t.Fatal("Cant connect", err)
	}

	contents, err := transport.ReadFile("new-os-release")
	if err != nil {
		t.Fatal(err)
	}

	strings.Contains(string(contents[:]), "SSH Linux")
}
