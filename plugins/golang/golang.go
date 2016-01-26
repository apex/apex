// Package golang implements the "golang" runtime.
package golang

import (
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("golang", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Run adds the shim when runtime is "golang".
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if hook != function.OpenHook || fn.Runtime != "golang" {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "GOOS=linux GOARCH=amd64 go build -o main main.go"
	}

	fn.Shim = true
	fn.Runtime = "nodejs"
	fn.Hooks.Clean = "rm -f main"

	return nil
}
