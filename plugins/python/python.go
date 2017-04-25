// Package python implements the "python" runtime.
package python

import (
	"strings"

	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("python", &Plugin{})
}

// Runtime for inference.
const Runtime = "python3.6"

// Plugin implementation.
type Plugin struct{}

// Open adds python defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "python") {
		return nil
	}

	// Support "python" for backwards compat.
	if fn.Runtime == "python" {
		fn.Runtime = "python2.7"
	}

	if fn.Handler == "" {
		fn.Handler = "main.handle"
	}

	return nil
}
