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

// Open adds nodejs defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != "nodejs" {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}
