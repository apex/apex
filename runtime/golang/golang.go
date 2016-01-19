package golang

import (
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
	cmd := exec.Command("go", "build", "-o", "main", "main.go")
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}

func (r *Runtime) Clean(dir string) error {
	return os.Remove(filepath.Join(dir, "main"))
}

func (r *Runtime) DefaultFile() string {
	return "main.go"
}
