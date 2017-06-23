// Package nodejs implements the "nodejs" runtime.
package nodejs

import (
	"github.com/apex/apex/function"
	"strings"
)

const (
	// Runtime name used by Apex and by AWS Lambda for Node.js 0.10
	Runtime = "nodejs"
)

func init() {
	function.RegisterPlugin(Runtime, &Plugin{})
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
