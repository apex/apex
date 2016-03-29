// Package java implements the "java" runtime.
package java

import (
	azip "archive/zip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/apex/function"
	"github.com/jpillora/archive"
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
	if fn.Runtime != Runtime {
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
func (p *Plugin) Build(fn *function.Function, zip *archive.Archive) error {
	if fn.Runtime != Runtime {
		return nil
	}
	fn.Runtime = RuntimeCanonical

	fn.Log.Debugf("searching for JAR (%s) in directories: %s", jarFile, strings.Join(jarSearchPaths, ", "))
	expectedJarPath := findJar(fn.Path)
	if expectedJarPath == "" {
		return errors.New("Expected jar file not found")
	}
	fn.Log.Debugf("found jar path: %s", expectedJarPath)

	fn.Log.Debug("appending compiled files")
	reader, err := azip.OpenReader(expectedJarPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		r, err := file.Open()
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		r.Close()

		zip.AddBytes(file.Name, b)
	}

	return nil
}

func findJar(fnPath string) string {
	for _, path := range jarSearchPaths {
		jarPath := filepath.Join(fnPath, path, jarFile)
		if _, err := os.Stat(jarPath); err == nil {
			return jarPath
		}
	}
	return ""
}
