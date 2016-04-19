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

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/tj/go-sync/semaphore"
	"gopkg.in/validator.v2"

	"github.com/apex/apex/function"
	"github.com/apex/apex/hooks"
	"github.com/apex/apex/infra"
	"github.com/apex/apex/utils"
	"github.com/apex/apex/vpc"
)

const (
	// DefaultMemory defines default memory value (MB) for every function in a project
	DefaultMemory = 128

	// DefaultTimeout defines default timeout value (s) for every function in a project
	DefaultTimeout = 3
)

// Config for project.
type Config struct {
	Name               string            `json:"name" validate:"nonzero"`
	Description        string            `json:"description"`
	Runtime            string            `json:"runtime"`
	Memory             int64             `json:"memory"`
	Timeout            int64             `json:"timeout"`
	Role               string            `json:"role"`
	Handler            string            `json:"handler"`
	Shim               bool              `json:"shim"`
	NameTemplate       string            `json:"nameTemplate"`
	RetainedVersions   int               `json:"retainedVersions"`
	DefaultEnvironment string            `json:"defaultEnvironment"`
	Environment        map[string]string `json:"environment"`
	Hooks              hooks.Hooks       `json:"hooks"`
	VPC                vpc.VPC           `json:"vpc"`
}

// Project represents zero or more Lambda functions.
type Project struct {
	Config
	Path         string
	Alias        string
	Concurrency  int
	Environment  string
	Log          log.Interface
	Service      lambdaiface.LambdaAPI
	Functions    []*function.Function
	IgnoreFile   []byte
	nameTemplate *template.Template
}

// defaults applies configuration defaults.
func (p *Project) defaults() {
	p.Memory = DefaultMemory
	p.Timeout = DefaultTimeout
	p.IgnoreFile = []byte(".apexignore\nfunction.json\n")

	if p.Concurrency == 0 {
		p.Concurrency = 5
	}

	if p.Config.Environment == nil {
		p.Config.Environment = make(map[string]string)
	}

	if p.NameTemplate == "" {
		p.NameTemplate = "{{.Project.Name}}_{{.Function.Name}}"
	}

	if p.RetainedVersions == 0 {
		p.RetainedVersions = 10
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

	if p.Environment == "" {
		p.Environment = p.Config.DefaultEnvironment
	}
	if p.Role == "" {
		p.Role = p.readInfraRole()
	}

	if err := validator.Validate(&p.Config); err != nil {
		return err
	}

	t, err := template.New("nameTemplate").Parse(p.NameTemplate)
	if err != nil {
		return err
	}
	p.nameTemplate = t

	ignoreFile, err := utils.ReadIgnoreFile(p.Path)
	if err != nil {
		return err
	}
	p.IgnoreFile = append(p.IgnoreFile, ignoreFile...)

	return nil
}

// LoadFunctions reads the ./functions directory, populating the Functions field.
// If no `names` are specified, all functions are loaded.
func (p *Project) LoadFunctions(names ...string) error {
	dir := filepath.Join(p.Path, functionsDir)
	p.Log.Debugf("loading functions in %s", dir)

	existing, err := p.FunctionDirNames()
	if err != nil {
		return err
	}

	if len(names) == 0 {
		names = existing
	}

	for _, name := range names {
		if !utils.ContainsString(existing, name) {
			p.Log.Warnf("function %q does not exist in project", name)
			continue
		}

		fn, err := p.LoadFunction(name)
		if err != nil {
			return err
		}

		p.Functions = append(p.Functions, fn)
	}

	if len(p.Functions) == 0 {
		return errors.New("no function loaded")
	}

	return nil
}

// DeployAndClean deploys functions and then cleans up their build artifacts.
func (p *Project) DeployAndClean() error {
	if err := p.Deploy(); err != nil {
		return err
	}

	return p.Clean()
}

// Deploy functions and their configurations.
func (p *Project) Deploy() error {
	p.Log.Debugf("deploying %d functions", len(p.Functions))

	sem := make(semaphore.Semaphore, p.Concurrency)
	errs := make(chan error)

	go func() {
		for _, fn := range p.Functions {
			fn := fn
			sem.Acquire()

			go func() {
				defer sem.Release()

				err := fn.Deploy()
				if err != nil {
					err = fmt.Errorf("function %s: %s", fn.Name, err)
				}

				errs <- err
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

// Clean up function build artifacts.
func (p *Project) Clean() error {
	p.Log.Debugf("cleaning %d functions", len(p.Functions))

	for _, fn := range p.Functions {
		if err := fn.Clean(); err != nil {
			return fmt.Errorf("function %s: %s", fn.Name, err)
		}
	}

	return nil
}

// Delete functions.
func (p *Project) Delete() error {
	p.Log.Debugf("deleting %d functions", len(p.Functions))

	for _, fn := range p.Functions {
		if _, err := fn.GetConfig(); err != nil {
			if awserr, ok := err.(awserr.Error); ok && awserr.Code() == "ResourceNotFoundException" {
				p.Log.Infof("function %q hasn't been deployed yet or has been deleted manually on AWS Lambda", fn.Name)
				continue
			}
			return fmt.Errorf("function %s: %s", fn.Name, err)
		}

		if err := fn.Delete(); err != nil {
			return fmt.Errorf("function %s: %s", fn.Name, err)
		}
	}

	return nil
}

// Rollback project functions to previous version.
func (p *Project) Rollback() error {
	p.Log.Debugf("rolling back %d functions", len(p.Functions))

	for _, fn := range p.Functions {
		if err := fn.Rollback(); err != nil {
			return fmt.Errorf("function %s: %s", fn.Name, err)
		}
	}

	return nil
}

// RollbackVersion project functions to the specified version.
func (p *Project) RollbackVersion(version string) error {
	p.Log.Debugf("rolling back %d functions to version %s", len(p.Functions), version)

	for _, fn := range p.Functions {
		if err := fn.RollbackVersion(version); err != nil {
			return fmt.Errorf("function %s: %s", fn.Name, err)
		}
	}

	return nil
}

// FunctionDirNames returns a list of function directory names.
func (p *Project) FunctionDirNames() (list []string, err error) {
	dir := filepath.Join(p.Path, functionsDir)

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

// Setenv sets environment variable `name` to `value` on every function in project.
func (p *Project) Setenv(name, value string) {
	for _, fn := range p.Functions {
		fn.Setenv(name, value)
	}
}

// LoadFunction returns the function in the ./functions/<name> directory.
func (p *Project) LoadFunction(name string) (*function.Function, error) {
	return p.LoadFunctionByPath(name, filepath.Join(p.Path, functionsDir, name))
}

// LoadFunctionByPath returns the function in the given directory.
func (p *Project) LoadFunctionByPath(name, path string) (*function.Function, error) {
	p.Log.Debugf("loading function in %s", path)

	fn := &function.Function{
		Config: function.Config{
			Runtime:          p.Runtime,
			Memory:           p.Memory,
			Timeout:          p.Timeout,
			Role:             p.Role,
			Handler:          p.Handler,
			Shim:             p.Shim,
			Hooks:            p.Hooks,
			Environment:      copyStringMap(p.Config.Environment),
			RetainedVersions: p.RetainedVersions,
			VPC:              p.VPC,
		},
		Name:       name,
		Path:       path,
		Service:    p.Service,
		Log:        p.Log,
		IgnoreFile: p.IgnoreFile,
		Alias:      p.Alias,
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

// readInfraRole reads lambda function IAM role from infrastructure
func (p *Project) readInfraRole() string {
	role, err := infra.Output(p.Environment, "lambda_function_role_id")
	if err != nil {
		p.Log.Debugf("couldn't read role from infrastructure: %s", err)
		return ""
	}

	return role
}

const functionsDir = "functions"

// render returns a string by executing template `t` against the given value `v`.
func render(t *template.Template, v interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, v); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// copyStringMap returns a copy of `in`.
func copyStringMap(in map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = v
	}
	return out
}
