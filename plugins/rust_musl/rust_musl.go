// Package rust_musl implements the "rust" runtime with musl dependencies.
package rust_musl

import (
	"fmt"
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/nodejs"
)

func init() {
	function.RegisterPlugin("rust-musl", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "rust-musl"
)

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and golang defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = fmt.Sprintf("cargo build --target=x86_64-unknown-linux-musl --release && mv target/x86_64-unknown-linux-musl/release/%v ./main", fn.Name)
	}

	fn.Shim = true
	fn.Runtime = nodejs.Runtime43

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
