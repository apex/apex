// Package project implements multi-function operations.
package project

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/apex/function"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// ErrNotFound is returned when a function cannot be found.
var ErrNotFound = errors.New("project: no function found")

// Config for project.
type Config struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Project represents zero or more Lambda functions.
type Project struct {
	Config
	Path        string
	Concurrency int
	Log         log.Interface
	Service     lambdaiface.LambdaAPI
	Functions   []*function.Function
}

// Open the project.json file and prime the config.
func (p *Project) Open() error {
	if p.Concurrency == 0 {
		p.Concurrency = 5
	}

	f, err := os.Open(filepath.Join(p.Path, "project.json"))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(f).Decode(&p.Config); err != nil {
		return err
	}

	return p.loadFunctions()
}

// DeployAndClean deploys functions and then cleans up their build artifacts.
func (p *Project) DeployAndClean(names []string) error {
	if err := p.Deploy(names); err != nil {
		return err
	}

	return p.Clean(names)
}

// Deploy functions and their configurations.
func (p *Project) Deploy(names []string) error {
	p.Log.Debugf("deploying %d functions", len(names))

	for _, name := range names {
		fn, err := p.FunctionByName(name)

		if err == ErrNotFound {
			continue
		}

		if err := fn.Deploy(); err != nil {
			return err
		}

		if err := fn.DeployConfig(); err != nil {
			return err
		}
	}

	return nil
}

// Clean up function build artifacts.
func (p *Project) Clean(names []string) error {
	p.Log.Debugf("cleaning %d functions", len(names))

	for _, name := range names {
		fn, err := p.FunctionByName(name)

		if err == ErrNotFound {
			continue
		}

		if err := fn.Clean(); err != nil {
			return err
		}
	}

	return nil
}

// Delete functions.
func (p *Project) Delete(names []string) error {
	p.Log.Debugf("deleting %d functions", len(names))

	for _, name := range names {
		fn, err := p.FunctionByName(name)

		if err == ErrNotFound {
			p.Log.Warnf("function %q does not exist", name)
			continue
		}

		if err := fn.Delete(); err != nil {
			return err
		}
	}

	return nil
}

// FunctionByName returns a function by `name` or returns ErrNotFound.
func (p *Project) FunctionByName(name string) (*function.Function, error) {
	for _, fn := range p.Functions {
		if fn.Name == name {
			return fn, nil
		}
	}

	return nil, ErrNotFound
}

// FunctionNames returns a list of function names sans-directory.
func (p *Project) FunctionNames() (list []string, err error) {
	dir := filepath.Join(p.Path, "functions")

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			list = append(list, file.Name())
		}
	}

	return list, nil
}

// loadFunctions reads the ./functions directory, populating the Functions field.
func (p *Project) loadFunctions() error {
	dir := filepath.Join(p.Path, "functions")
	p.Log.Debugf("loading functions in %s", dir)

	names, err := p.FunctionNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		fn, err := p.loadFunction(name)
		if err != nil {
			return err
		}

		p.Functions = append(p.Functions, fn)
	}

	return nil
}

// loadFunction returns the function in the ./functions/<name> directory.
func (p *Project) loadFunction(name string) (*function.Function, error) {
	dir := filepath.Join(p.Path, "functions", name)
	p.Log.Debugf("loading function %s", dir)

	fn := &function.Function{
		Path:    dir,
		Service: p.Service,
		Log:     p.Log.WithField("function", name),
	}

	if err := fn.Open(); err != nil {
		return nil, err
	}

	return fn, nil
}
