package transport

type Transport interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, contents []byte) error
	Run(command string) ([]byte, error)
}
