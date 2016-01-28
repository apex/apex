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

// Open adds python defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != "python" {
		return nil
	}

	fn.Runtime = "python2.7"

	if fn.Handler == "" {
		fn.Handler = "main.handle"
	}

	return nil
}
