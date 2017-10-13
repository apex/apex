// Package java implements the "java" runtime.
package java

import (
	azip "archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
	"github.com/pkg/errors"
)

const (
	// Runtime name used by Apex
	Runtime = "java"

	// RuntimeCanonical represents names used by AWS Lambda
	RuntimeCanonical = "java8"

	// jarFile
	jarFile = "apex.jar"
)

var jarSearchPaths = []string{
	"target",
	"build/libs",
}

func init() {
	function.RegisterPlugin(Runtime, &Plugin{})
}

// Plugin implementation
type Plugin struct{}

// Open adds java defaults. No clean operation is implemented, as it is
// assumed that the build tool generating the fat JAR will handle that workflow
// on its own.
func (p *Plugin) Open(fn *function.Function) error {
	if !strings.HasPrefix(fn.Runtime, "java") {
		return nil
	}

	if fn.Handler == "" {
		fn.Handler = "lambda.Main::handler"
	}

	if len(fn.IgnoreFile) == 0 {
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
	if !strings.HasPrefix(fn.Runtime, "java") {
		return nil
	}

	fn.Log.Debugf("searching for JAR (%s) in directories: %s", jarFile, strings.Join(jarSearchPaths, ", "))
	jar := findJar(fn.Path)
	if jar == "" {
		return errors.Errorf("missing jar file %q", jar)
	}

	fn.Log.Debug("appending compiled files")
	reader, err := azip.OpenReader(jar)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
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
	if !strings.HasPrefix(fn.Runtime, "java") {
		return nil
	}
	fn.Runtime = RuntimeCanonical

	return nil
}

func findJar(fnPath string) string {
	for _, path := range jarSearchPaths {
		jar := filepath.Join(fnPath, path, jarFile)
		if _, err := os.Stat(jar); err == nil {
			return jar
		}
	}

	return ""
}
