package clojure

import (
	azip "archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("clojure", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "clojure"

	// RuntimeCanonical represents names used by AWS Lambda
	RuntimeCanonical = "java8"

	jarFile = "apex.jar"
)

// Plugin Does plugin things
type Plugin struct{}

// Open adds the shim and golang defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "clojure") {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "lein uberjar && mv target/*-standalone.jar target/apex.jar"
	}

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "rm -fr target"
	}

	if _, err := os.Stat(".apexignore"); err != nil {
		// Since we're deploying a fat jar, we don't need anything else.
		fn.IgnoreFile = []byte(`
*
!**/apex.jar
`)
	}

	return nil
}

// Build adds the jar contents to zipfile.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	if !strings.HasPrefix(fn.Runtime, "clojure") {
		return nil
	}

	jar := filepath.Join(fn.Path, "target", jarFile)
	if _, err := os.Stat(jar); err != nil {
		return errors.Errorf("missing jar file %q", jar)
	}

	fn.Log.Debug("appending compiled files")
	reader, err := azip.OpenReader(jar)
	if err != nil {
		return errors.Wrap(err, "opening zip")
	}
	defer reader.Close()

	for _, file := range reader.File {
		parts := strings.Split(file.Name, ".")
		ext := parts[len(parts)-1]

		if ext == "clj" || ext == "cljx" || ext == "cljc" {
			continue
		}

		r, err := file.Open()
		if err != nil {
			return errors.Wrap(err, "opening file")
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			return errors.Wrap(err, "reading file")
		}
		r.Close()

		zip.AddBytes(file.Name, b)
	}

	return nil
}

func (p *Plugin) Deploy(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "clojure") {
		return nil
	}
	fn.Runtime = RuntimeCanonical

	return nil
}
