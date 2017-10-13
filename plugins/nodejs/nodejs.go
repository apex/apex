// Package nodejs implements the "nodejs" runtime.
package nodejs

import (
	"strings"

	"github.com/apex/apex/function"
)

const (
	// Runtime for inference.
	Runtime = "nodejs6.10"
)

func init() {
	function.RegisterPlugin("nodejs", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Open adds nodejs defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "nodejs") {
		return nil
	}

	if fn.Runtime == "nodejs" {
		fn.Runtime = Runtime
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}
