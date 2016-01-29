// Package shim adds a nodejs shim when the shim field is true, this is also used transparently
// in other plugins such as "golang" which are not directly supported by Lambda.
package shim

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/apex/function"
	"github.com/apex/apex/shim"
	"github.com/jpillora/archive"
)

func init() {
	function.RegisterPlugin("shim", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Build adds the nodejs shim files.
func (p *Plugin) Build(fn *function.Function, zip *archive.Archive) error {
	if fn.Shim {
		return p.addShim(fn)
	}

	return nil
}

// Clean removes the nodejs shim files.
func (p *Plugin) Clean(fn *function.Function) error {
	if fn.Shim {
		return p.removeShim(fn)
	}

	return nil
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
