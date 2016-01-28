// Package env populates .env.json if the function has any environment variables defined.
package env

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("env", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Build hook adds .env.json populate with Function.Enironment.
func (p *Plugin) Build(fn *function.Function) error {
	if len(fn.Environment) == 0 {
		return nil
	}

	return p.addEnv(filepath.Join(fn.Path, ".env.json"), fn)
}

// Clean hook removes .env.json.
func (p *Plugin) Clean(fn *function.Function) error {
	if len(fn.Environment) == 0 {
		return nil
	}

	return os.Remove(filepath.Join(fn.Path, ".env.json"))
}

// addEnv saves the environment as json into .env.json.
func (p *Plugin) addEnv(path string, fn *function.Function) error {
	fn.Log.WithField("env", fn.Environment).Debug("adding env")

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(fn.Environment)
}
