package golang

import (
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

func (r *Runtime) Compile() error {
	cmd := exec.Command("sh", "-c", `GOOS=linux GOARCH=amd64 go build -o main main.go`)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
