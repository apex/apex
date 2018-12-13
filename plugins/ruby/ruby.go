// Package ruby implements the "ruby" runtime.
package ruby

import (
	"strings"

	"github.com/apex/apex/function"
)

const (
	// Runtime for inference.
	Runtime = "ruby2.5"
)

func init() {
	function.RegisterPlugin("ruby", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Open adds ruby defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "ruby") {
		return nil
	}

	if fn.Runtime == "ruby" {
		fn.Runtime = Runtime
	}

	if fn.Handler == "" {
		fn.Handler = "lambda.handler"
	}

	return nil
}
