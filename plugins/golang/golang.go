// Package golang implements the "golang" runtime.
package golang

import (
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/nodejs"
)

func init() {
	function.RegisterPlugin("golang", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "golang"
)

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and golang defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "GOOS=linux GOARCH=amd64 go build -o main *.go"
	}

	fn.Shim = true
	fn.Runtime = nodejs.Runtime

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "rm -f main"
	}

	return nil
}
