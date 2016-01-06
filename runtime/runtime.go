// Package runtime provides interfaces for defining Lambda runtimes
// and appropriate shims for arbitrary language support.
package runtime

import (
	"errors"
)

// Registered runtimes.
var runtimes = make(map[string]Runtime)

// Runtime is a language runtime.
type Runtime interface {
	// Name returns the canonical runtime to be used, for example
	// since Go must be run as a shim, this is "nodejs", not "golang".
	Name() string

	// Handler returns the handler name for the runtime in the form "<file>.<func>".
	Handler() string

	// Shimmed returns true if the program should be shimmed.
	Shimmed() bool
}

// CompiledRuntime is a language runtime requiring compilation.
type CompiledRuntime interface {
	// Build performs a buid, using language-specific conventions
	// such as "main.go" for the build target.
	Build(dir string) error

	// Clean removes the build artifacts after deployment.
	Clean(dir string) error
}

// Register runtime by `name`.
func Register(name string, runtime Runtime) {
	runtimes[name] = runtime
}

// ByName returns the runtime by `name`.
func ByName(name string) (Runtime, error) {
	v, ok := runtimes[name]

	if !ok {
		return nil, errors.New("invalid runtime")
	}

	return v, nil
}
