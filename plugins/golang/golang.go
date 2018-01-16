// Package golang implements the "golang" runtime.
package golang

import (
	"strings"

	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("golang", &Plugin{})
}

const (
	// Runtime for inference.
	Runtime = "go1.x"
)

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and golang defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "go") {
		return nil
	}

	if fn.Runtime == "golang" {
		fn.Runtime = Runtime
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "GOOS=linux GOARCH=amd64 go build -o main *.go"
	}

	if fn.Handler == "" {
		fn.Handler = "main"
	}

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "rm -f main"
	}

	return nil
}
