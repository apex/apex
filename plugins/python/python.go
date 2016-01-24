// Package python implements the "python" runtime.
package python

import (
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("python", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Run specifies python defaults.
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if hook != function.OpenHook || fn.Runtime != "python" {
		return nil
	}

	fn.Runtime = "python2.7"

	if fn.Handler == "" {
		fn.Handler = "main.handle"
	}

	return nil
}
