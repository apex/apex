// Package env populates .env.json if the function has any environment variables defined.
package env

import (
	"encoding/json"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("env", &Plugin{})
}

// FileName of file with environment variables.
const FileName = ".env.json"

// Plugin implementation.
type Plugin struct{}

// Build hook adds .env.json populate with Function.Enironment.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	if len(fn.Environment) == 0 {
		return nil
	}

	fn.Log.WithField("env", fn.Environment).Debug("adding env")

	env, err := json.Marshal(fn.Environment)
	if err != nil {
		return err
	}

	return zip.AddBytes(FileName, env)
}
