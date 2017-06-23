// Package rust_gnu implements the "rust" runtime with gnu libc dependencies.
package rust_gnu

import (
	"fmt"
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/nodejs"
	"strings"
)

func init() {
	function.RegisterPlugin("rust-gnu", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "rust-gnu"
)

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and golang defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "rust-gnu") {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = fmt.Sprintf("cargo build --target=x86_64-unknown-linux-gnu --release && mv target/x86_64-unknown-linux-gnu/release/%v ./main", fn.Name)
	}

	fn.Shim = true
	fn.Runtime = nodejs.Runtime

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "rm -f main"
	}

	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\ntarget/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nsrc/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nexamples/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\ntests/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nbenches/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nCargo.toml")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nCargo.lock")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\n.git/")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\n.gitignore")...)
	fn.IgnoreFile = append(fn.IgnoreFile, []byte("\nbuild.rs")...)

	return nil
}
