package transport

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshTransport struct {
	config     *ssh.ClientConfig
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewSshTransport(host string, config *ssh.ClientConfig) (*SshTransport, error) {
	sshClient, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	// TODO
	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return &SshTransport{
		config:     config,
		sshClient:  sshClient,
		sftpClient: sftpClient,
	}, nil
}

func (s *SshTransport) session() *ssh.Session {
	session, err := s.sshClient.NewSession()
	if err != nil {
		panic(err)
	}
	return session
}

func (s *SshTransport) ReadFile(path string) ([]byte, error) {
	file, err := s.sftpClient.Open(path)
	if err != nil {
		return nil, err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func (s *SshTransport) WriteFile(path string, contents []byte) error {
	file, err := s.sftpClient.OpenFile(path, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return err
	}

	// TODO write to temp file + move once verified

	bytesWritten, err := file.ReadFrom(bytes.NewReader(contents))
	if err != nil {
		return err
	}

	if bytesWritten != int64(len(contents)) {
		return errors.New("File write incomplete")
	}
	return nil
}

func (s *SshTransport) Run(command string) ([]byte, error) {
	return s.session().Output(command)
}
