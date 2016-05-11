// Package shim adds a nodejs shim when the shim field is true, this is also used transparently
// in other plugins such as "golang" which are not directly supported by Lambda.
package shim

import (
	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
	"github.com/apex/apex/shim"
)

func init() {
	function.RegisterPlugin("shim", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Build adds the nodejs shim files.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	if fn.Shim {
		fn.Log.Debug("add shim")

		if err := zip.AddBytes("index.js", shim.MustAsset("index.js")); err != nil {
			return err
		}

		if err := zip.AddBytes("byline.js", shim.MustAsset("byline.js")); err != nil {
			return err
		}
	}

	return nil
}
