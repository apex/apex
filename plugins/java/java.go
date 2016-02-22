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
	err := mvnCmd.Run()
	if err != nil {
		return err
	}

	expectedJarPath := filepath.Join(fn.Path, "target", targetJarFile+".jar")
	if _, err := os.Stat(expectedJarPath); err != nil {
		return errors.New("Expected jar file not found")
	}

	fn.Log.Debug("appending compiled files")
	reader, err := azip.OpenReader(expectedJarPath)
	defer reader.Close()
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		fileContents, err := ioutil.ReadAll(fileReader)
		fileReader.Close()
		if err != nil {
			return err
		}
		zip.AddBytes(file.Name, fileContents)
	}

	fn.Log.Debug("cleaning mvn tmpfiles")
	mvnCleanCmd := exec.Command("mvn", "clean")
	mvnCleanCmd.Dir = fn.Path
	err = mvnCleanCmd.Run()
	if err != nil {
		return err
	}

	if generatedPom {
		os.Remove(expectedPomPath)
	}

	return nil
}
