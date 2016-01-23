// Package project implements multi-function operations.
package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/tj/go-sync/semaphore"
	"gopkg.in/validator.v2"

	"github.com/apex/apex/function"
	"github.com/apex/apex/hooks"
	"github.com/apex/apex/logs"
	"github.com/apex/apex/utils"
	"github.com/apex/log"
)

const (
	// DefaultMemory defines default memory value (MB) for every function in a project
	DefaultMemory = 128

	// DefaultTimeout defines default timeout value (s) for every function in a project
	DefaultTimeout = 3
)

// ErrNotFound is returned when a function cannot be found.
var ErrNotFound = errors.New("project: no function found")

// Config for project.
type Config struct {
	Name         string      `json:"name" validate:"nonzero"`
	Description  string      `json:"description"`
	Runtime      string      `json:"runtime"`
	Memory       int64       `json:"memory"`
	Timeout      int64       `json:"timeout"`
	Role         string      `json:"role"`
	NameTemplate string      `json:"nameTemplate"`
	Hooks        hooks.Hooks `json:"hooks"`
}

// Project represents zero or more Lambda functions.
type Project struct {
	Config
	Path            string
	Concurrency     int
	Log             log.Interface
	Service         lambdaiface.LambdaAPI
	Functions       []*function.Function
	IgnoredPatterns []string
	nameTemplate    *template.Template
}

// defaults applies configuration defaults.
func (p *Project) defaults() {
	p.Memory = DefaultMemory
	p.Timeout = DefaultTimeout

	if p.Concurrency == 0 {
		p.Concurrency = 5
	}

	if p.NameTemplate == "" {
		p.NameTemplate = "{{.Project.Name}}_{{.Function.Name}}"
	}
}

// Open the project.json file and prime the config.
func (p *Project) Open() error {
	p.defaults()

	f, err := os.Open(filepath.Join(p.Path, "project.json"))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(f).Decode(&p.Config); err != nil {
		return err
	}

	if err := validator.Validate(&p.Config); err != nil {
		return err
	}

	t, err := template.New("nameTemplate").Parse(p.NameTemplate)
	if err != nil {
		return err
	}
	p.nameTemplate = t

	p.IgnoredPatterns, err = utils.ReadIgnoreFile(p.Path)
	if err != nil {
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

	sem := make(semaphore.Semaphore, p.Concurrency)
	errs := make(chan error)

	go func() {
		for _, name := range names {
			name := name
			sem.Acquire()

			go func() {
				defer sem.Release()
				errs <- p.deploy(name)
			}()
		}

		sem.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// deploy function by `name`.
func (p *Project) deploy(name string) error {
	fn, err := p.FunctionByName(name)

	if err == ErrNotFound {
		p.Log.Warnf("function %q does not exist", name)
		return nil
	}

	if err != nil {
		return err
	}

	return fn.Deploy()
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
			p.Log.Warnf("function %q does not exist in project", name)
			continue
		}

		if _, err := fn.GetConfig(); err != nil {
			if awserr, ok := err.(awserr.Error); ok && awserr.Code() == "ResourceNotFoundException" {
				p.Log.Infof("function %q hasn't been deployed yet or has been deleted manually on AWS Lambda", name)
				continue
			}
			return err
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

// FunctionDirNames returns a list of function directory names.
func (p *Project) FunctionDirNames() (list []string, err error) {
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

// FunctionNames returns a list of function names.
func (p *Project) FunctionNames() (list []string) {
	for _, fn := range p.Functions {
		list = append(list, fn.Name)
	}

	return list
}

// Logs returns logs.
func (p *Project) Logs(s *session.Session, name string, filter string) (*logs.Logs, error) {
	fn, err := p.FunctionByName(name)
	if err != nil {
		return nil, err
	}
	fnName, err := p.name(fn)
	if err != nil {
		return nil, err
	}
	l := &logs.Logs{
		Service:       cloudwatchlogs.New(s),
		Log:           log.Log,
		GroupName:     fmt.Sprintf("/aws/lambda/%s", fnName),
		FilterPattern: filter,
	}
	return l, nil
}

// SetEnv sets environment variable `name` to `value` on every function in project.
func (p *Project) SetEnv(name, value string) {
	for _, fn := range p.Functions {
		fn.SetEnv(name, value)
	}
}

// loadFunctions reads the ./functions directory, populating the Functions field.
func (p *Project) loadFunctions() error {
	dir := filepath.Join(p.Path, "functions")
	p.Log.Debugf("loading functions in %s", dir)

	names, err := p.FunctionDirNames()
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
	p.Log.Debugf("loading function in %s", dir)

	fn := &function.Function{
		Config: function.Config{
			Runtime: p.Runtime,
			Memory:  p.Memory,
			Timeout: p.Timeout,
			Role:    p.Role,
			Hooks:   p.Hooks,
		},
		Name:            name,
		Path:            dir,
		Service:         p.Service,
		Log:             p.Log,
		IgnoredPatterns: p.IgnoredPatterns,
	}

	if name, err := p.name(fn); err == nil {
		fn.FunctionName = name
	} else {
		return nil, err
	}

	if err := fn.Open(); err != nil {
		return nil, err
	}

	return fn, nil
}

// name returns the computed name for `fn`, using the nameTemplate.
func (p *Project) name(fn *function.Function) (string, error) {
	data := struct {
		Project  *Project
		Function *function.Function
	}{
		Project:  p,
		Function: fn,
	}

	name, err := render(p.nameTemplate, data)
	if err != nil {
		return "", err
	}

	return name, nil
}

// render returns a string by executing template `t` against the given value `v`.
func render(t *template.Template, v interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, v); err != nil {
		return "", err
	}

	return buf.String(), nil
}
