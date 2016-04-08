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

const (
	// Runtime name used by Apex and by AWS Lambda for Node.js 0.10
	Runtime = "nodejs"
	// Runtime43 name used by Apex and by AWS Lambda for Node.js 4.3.2
	Runtime43 = "nodejs4.3"
)

func init() {
	function.RegisterPlugin(Runtime, &Plugin{})
	function.RegisterPlugin(Runtime43, &Plugin{})
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

// Plugin implementation.
type Plugin struct{}

// Open adds nodejs defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !p.runtimeSupported(fn) {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "index.handle"
	}

	return nil
}

// Build injects a script for loading the environment.
func (p *Plugin) Build(fn *function.Function, zip *archive.Archive) error {
	if !p.runtimeSupported(fn) || len(fn.Environment) == 0 {
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

func (p *Plugin) runtimeSupported(fn *function.Function) bool {
	return fn.Runtime == Runtime || fn.Runtime == Runtime43
}
