package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/apex/apex/runtime"
)

func init() {
	runtime.Register("golang", new(Runtime))
}

type Runtime struct{}

func (r *Runtime) Name() string {
	return "nodejs"
}

func (r *Runtime) Shimmed() bool {
	return true
}

func (r *Runtime) Handler() string {
	return "index.handle"
}

func (r *Runtime) Build(dir string) error {
	s := fmt.Sprintf("cd %s && GOOS=linux GOARCH=amd64 go build -o main main.go", dir)
	cmd := exec.Command("sh", "-c", s)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runtime) Clean(dir string) error {
	return os.Remove(filepath.Join(dir, "main"))
}
