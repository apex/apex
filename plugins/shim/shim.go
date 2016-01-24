// Package shim adds a nodejs shim when the shim field is true, this is also used transparently
// in other plugins such as "golang" which are not directly supported by Lambda.
package shim

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/apex/function"
	"github.com/apex/apex/shim"
)

func init() {
	function.RegisterPlugin("shim", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Run adds the shim on build and removes it on clean.
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if !fn.Shim {
		return nil
	}

	switch hook {
	case function.BuildHook:
		return p.addShim(fn)
	case function.CleanHook:
		return p.removeShim(fn)
	default:
		return nil
	}
}

// addShim saves nodejs shim.
func (p *Plugin) addShim(fn *function.Function) error {
	fn.Log.Debug("add shim")

	if err := ioutil.WriteFile(filepath.Join(fn.Path, "index.js"), shim.MustAsset("index.js"), 0666); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(fn.Path, "byline.js"), shim.MustAsset("byline.js"), 0666); err != nil {
		return err
	}

	return nil
}

// removeShim removes the nodejs shim.
func (p *Plugin) removeShim(fn *function.Function) error {
	if err := os.Remove(filepath.Join(fn.Path, "index.js")); err != nil {
		return err
	}

	if err := os.Remove(filepath.Join(fn.Path, "byline.js")); err != nil {
		return err
	}

	return nil
}
