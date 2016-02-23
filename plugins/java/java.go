// Package java implements the "java" runtime.
package java

import (
	azip "archive/zip"

	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"errors"

	"github.com/apex/apex/function"
	"github.com/jpillora/archive"
)

const (
	// Runtime name used by Apex
	Runtime = "java"
	// RuntimeCanonical represents names used by AWS Lambda
	RuntimeCanonical = "java8"
	// targetJarFile mvn target jar file name
	targetJarFile = "apex-plugin-target"
)

func init() {
	function.RegisterPlugin("java", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// Open adds java defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	fn.Runtime = RuntimeCanonical

	if fn.Handler == "" {
		fn.Handler = "lambda.Main::handler"
	}

	return nil
}

// Build calls mvn package, add jar contents to zipfile.
func (p *Plugin) Build(fn *function.Function, zip *archive.Archive) error {
	if fn.Runtime != RuntimeCanonical {
		return nil
	}

	generatedPom := false

	expectedPomPath := filepath.Join(fn.Path, "pom.xml")
	if _, err := os.Stat(expectedPomPath); err != nil {
		fn.Log.Debug("generating default pom")
		generatedPom = true
		if err := ioutil.WriteFile(expectedPomPath, []byte(genericPom), 0644); err != nil {
			return err
		}
	}

	fn.Log.Debug("creating jar")
	mvnCmd := exec.Command("mvn", "package", "-Djar.finalName="+targetJarFile)
	mvnCmd.Dir = fn.Path
	if err := mvnCmd.Run(); err != nil {
		return err
	}

	expectedJarPath := filepath.Join(fn.Path, "target", targetJarFile+".jar")
	if _, err := os.Stat(expectedJarPath); err != nil {
		return errors.New("Expected jar file not found")
	}

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

	if err := p.cleanTarget(fn); err != nil {
		return err
	}

	if generatedPom {
		os.Remove(expectedPomPath)
	}

	return nil
}

// cleanTarget cleans target dir. pom.xml is required to call this function.
func (p *Plugin) cleanTarget(fn *function.Function) error {
	fn.Log.Debug("cleaning mvn tmpfiles")
	cmd := exec.Command("mvn", "clean")
	cmd.Dir = fn.Path
	return cmd.Run()
}
