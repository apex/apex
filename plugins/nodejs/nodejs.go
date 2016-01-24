// Package nodejs implements the "nodejs" runtime.
package nodejs

import (
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("nodejs", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Run specifies nodejs defaults.
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if hook != function.OpenHook || fn.Runtime != "nodejs" {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}
