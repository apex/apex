// Package python implements the "python" runtime.
package python

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/env"
)

func init() {
	function.RegisterPlugin("python", &Plugin{})
}

// prelude script template.
var prelude = template.Must(template.New("prelude").Parse(`import json
import os

try:
    with open('./{{.EnvFile}}') as env_file:
        config = json.load(env_file)
    for key, value in config.items():
        os.environ[key] = str(value)
except IOError:
    pass

from {{.HandleFile}} import {{.HandleMethod}}
`))

const (
	// Runtime name used by Apex
	Runtime = "python"
	// RuntimeCanonical represents names used by AWS Lambda
	RuntimeCanonical = "python2.7"
)

// Plugin implementation.
type Plugin struct{}

// Open adds python defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	fn.Runtime = RuntimeCanonical

	if fn.Handler == "" {
		fn.Handler = "main.handle"
	}

	return nil
}

// Build injects a script for loading the environment.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	if fn.Runtime != RuntimeCanonical || len(fn.Environment) == 0 {
		return nil
	}

	fn.Log.Debug("injecting prelude")

	var buf bytes.Buffer
	file := strings.Split(fn.Handler, ".")[0]
	method := strings.Split(fn.Handler, ".")[1]

	err := prelude.Execute(&buf, struct {
		EnvFile      string
		HandleFile   string
		HandleMethod string
	}{
		EnvFile:      env.FileName,
		HandleFile:   file,
		HandleMethod: method,
	})

	if err != nil {
		return err
	}

	fn.Handler = "_apex_main." + method

	return zip.AddBytes("_apex_main.py", buf.Bytes())
}
