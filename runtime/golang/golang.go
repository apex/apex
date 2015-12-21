package golang

import (
	"fmt"
	"os"
	"os/exec"
)

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

func (r *Runtime) Build(target string) error {
	if target == "" {
		target = "main.go"
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`GOOS=linux GOARCH=amd64 go build -o main %s`, target))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runtime) Clean() error {
	return os.Remove("main")
}
