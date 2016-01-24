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

// Run adds .env.json on build and removes it on clean.
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if len(fn.Environment) == 0 {
		return nil
	}

	path := filepath.Join(fn.Path, ".env.json")

	switch hook {
	case function.BuildHook:
		return p.addEnv(path, fn)
	case function.CleanHook:
		return os.Remove(path)
	default:
		return nil
	}
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
