// Package shim adds a nodejs shim when the shim field is true, this is also used transparently
// in other plugins such as "golang" which are not directly supported by Lambda.
package shim

import (
	"time"

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
		fn.Log.Debug("add shim")

		if err := zip.AddBytesMTime("index.js", shim.MustAsset("index.js"), time.Unix(0, 0)); err != nil {
			return err
		}

		if err := zip.AddBytesMTime("byline.js", shim.MustAsset("byline.js"), time.Unix(0, 0)); err != nil {
			return err
		}
	}

	return nil
}
