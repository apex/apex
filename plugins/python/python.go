// Package python implements the "python" runtime.
package python

import "github.com/apex/apex/function"

func init() {
	function.RegisterPlugin("python", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "python"

	// RuntimeCanonical represents names used by AWS Lambda
	RuntimeCanonical = "python2.7"
)

// Plugin implementation.
type Plugin struct{}

// Open adds python defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	fn.Runtime = RuntimeCanonical

	if fn.Handler == "" {
		fn.Handler = "main.handle"
	}

	return nil
}
