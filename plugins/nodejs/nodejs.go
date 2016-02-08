// Package nodejs implements the "nodejs" runtime.
package nodejs

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/jpillora/archive"

	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/env"
)

func init() {
	function.RegisterPlugin("nodejs", &Plugin{})
}

// prelude script template.
var prelude = template.Must(template.New("prelude").Parse(`try {
  var config = require('./{{.EnvFile}}')
  for (var key in config) {
    process.env[key] = config[key]
  }
} catch (err) {
  // ignore
}

exports.handle = require('./{{.HandleFile}}').{{.HandleMethod}}
`))

const (
	// Runtime name used by Apex and by AWS Lambda
	Runtime = "nodejs"
)

// Plugin implementation.
type Plugin struct{}

// Open adds nodejs defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}

// Build injects a script for loading the environment.
func (p *Plugin) Build(fn *function.Function, zip *archive.Archive) error {
	if fn.Runtime != Runtime || len(fn.Environment) == 0 {
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

	fn.Handler = "_apex_index.handle"

	return zip.AddBytesMTime("_apex_index.js", buf.Bytes(), time.Unix(0, 0))
}
