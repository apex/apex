// Package nodejs implements the "nodejs" runtime.
package nodejs

import "github.com/apex/apex/function"

const (
	// Runtime name used by Apex and by AWS Lambda for Node.js 0.10
	Runtime = "nodejs"

	// Runtime43 name used by Apex and by AWS Lambda for Node.js 4.3.2
	Runtime43 = "nodejs4.3"

	// Runtime43Edge name used by Apex and by AWS Lambda for Node.js 4.3.2 Edge
	Runtime43Edge = "nodejs4.3-edge"

	// Runtime6_10 name used by Apex and by AWS Lambda for Node.js 6.10
	Runtime6_10 = "nodejs6.10"
)

func init() {
	function.RegisterPlugin(Runtime, &Plugin{})
	function.RegisterPlugin(Runtime43, &Plugin{})
	function.RegisterPlugin(Runtime43Edge, &Plugin{})
	function.RegisterPlugin(Runtime6_10, &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Open adds nodejs defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !runtimeSupported(fn) {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}

func runtimeSupported(fn *function.Function) bool {
	return fn.Runtime == Runtime || fn.Runtime == Runtime43 || fn.Runtime == Runtime43Edge || fn.Runtime == Runtime6_10
}
