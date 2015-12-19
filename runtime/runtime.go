// Package runtime provides interfaces for defining Lambda runtimes
// and appropriate shims for arbitrary language support.
package runtime

import (
	"errors"

	"github.com/apex/apex/runtime/golang"
	"github.com/apex/apex/runtime/nodejs"
)

// Runtime is a language runtime.
type Runtime interface {
	// Name returns the canonical runtime to be used, for example
	// since Go must be run as a shim, this is "nodejs", not "golang".
	Name() string

	// Shimmed returns true if the program should be shimmed.
	Shimmed() bool
}

// CompiledRuntime is a language runtime requiring compilation.
type CompiledRuntime interface {
	Compile() error
}

// runtimes map by name.
var runtimes = map[string]Runtime{
	"nodejs": &nodejs.Runtime{},
	"golang": &golang.Runtime{},
}

// ByName returns the runtime by `name`.
func ByName(name string) (Runtime, error) {
	v, ok := runtimes[name]

	if !ok {
		return nil, errors.New("invalid runtime")
	}

	return v, nil
}
