// Package function implements function-level operations.
package function

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/dustin/go-humanize"
	"github.com/jpillora/archive"
	"gopkg.in/validator.v2"

	"github.com/apex/apex/hooks"
	"github.com/apex/apex/runtime"
	"github.com/apex/apex/shim"
	"github.com/apex/apex/utils"
	"github.com/apex/log"
)

// InvocationType determines how an invocation request is made.
type InvocationType string

// Invocation types.
const (
	RequestResponse InvocationType = "RequestResponse"
	Event                          = "Event"
	DryRun                         = "DryRun"
)

// Current alias name.
const CurrentAlias = "current"

// InvokeError records an error from an invocation.
type InvokeError struct {
	Message string   `json:"errorMessage"`
	Type    string   `json:"errorType"`
	Stack   []string `json:"stackTrace"`
	Handled bool
}

// Error message.
func (e *InvokeError) Error() string {
	return e.Message
}

// Config for a Lambda function.
type Config struct {
	Description string      `json:"description"`
	Runtime     string      `json:"runtime" validate:"nonzero"`
	Memory      int64       `json:"memory" validate:"nonzero"`
	Timeout     int64       `json:"timeout" validate:"nonzero"`
	Role        string      `json:"role" validate:"nonzero"`
	Hooks       hooks.Hooks `json:"hooks"`
}

// Function represents a Lambda function, with configuration loaded
// from the "function.json" file on disk. Operations are performed
// against the function directory as the CWD, so os.Chdir() first.
type Function struct {
	Config
	Name            string
	FunctionName    string
	Path            string
	Service         lambdaiface.LambdaAPI
	Log             log.Interface
	IgnoredPatterns []string
	runtime         runtime.Runtime
	env             map[string]string
}

// Open the function.json file and prime the config.
func (f *Function) Open() error {
	p, err := os.Open(filepath.Join(f.Path, "function.json"))
	if err == nil {
		if err := json.NewDecoder(p).Decode(&f.Config); err != nil {
			return err
		}
	}

	if f.Runtime == "" {
		if runtimeName, err := runtime.Detect(f.Path); err == nil {
			f.Runtime = runtimeName
		}
	}

	if err := validator.Validate(&f.Config); err != nil {
		return fmt.Errorf("error opening function %s: %s", f.Name, err.Error())
	}

	r, err := runtime.ByName(f.Runtime)
	if err != nil {
		return err
	}
	f.runtime = r

	patterns, err := utils.ReadIgnoreFile(f.Path)
	if err != nil {
		return err
	}
	f.IgnoredPatterns = append(f.IgnoredPatterns, patterns...)

	f.Log = f.Log.WithField("function", f.Name)

	return nil
}

// SetEnv sets environment variable `name` to `value`.
func (f *Function) SetEnv(name, value string) {
	if f.env == nil {
		f.env = make(map[string]string)
	}
	f.env[name] = value
}

// Deploy code and then configuration.
func (f *Function) Deploy() error {
	if err := f.DeployCode(); err != nil {
		return err
	}

	return f.DeployConfig()
}

// DeployCode generates a zip and creates or updates the function.
func (f *Function) DeployCode() error {
	f.Log.Info("deploying")

	zip, err := f.BuildBytes()
	if err != nil {
		return err
	}

	config, err := f.GetConfig()

	if e, ok := err.(awserr.Error); ok {
		if e.Code() == "ResourceNotFoundException" {
			return f.Create(zip)
		}
	}

	if err != nil {
		return err
	}

	remoteHash := *config.Configuration.CodeSha256
	localHash := utils.Sha256(zip)

	if localHash == remoteHash {
		f.Log.Info("unchanged")
		return nil
	}

	return f.Update(zip)
}

// DeployConfig deploys changes to configuration.
func (f *Function) DeployConfig() error {
	f.Log.Info("deploying config")

	_, err := f.Service.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: &f.FunctionName,
		MemorySize:   &f.Memory,
		Timeout:      &f.Timeout,
		Description:  &f.Description,
		Role:         aws.String(f.Role),
		Handler:      aws.String(f.runtime.Handler()),
	})

	return err
}

// Delete the function including all its versions
func (f *Function) Delete() error {
	f.Log.Info("deleting")
	_, err := f.Service.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: &f.FunctionName,
	})
	return err
}

// GetConfig returns the function configuration.
func (f *Function) GetConfig() (*lambda.GetFunctionOutput, error) {
	f.Log.Debug("fetching config")
	return f.Service.GetFunction(&lambda.GetFunctionInput{
		FunctionName: &f.FunctionName,
	})
}

// Update the function with the given `zip`.
func (f *Function) Update(zip []byte) error {
	f.Log.Info("updating function")

	updated, err := f.Service.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: &f.FunctionName,
		Publish:      aws.Bool(true),
		ZipFile:      zip,
	})

	if err != nil {
		return err
	}

	f.Log.Info("updating alias")

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &f.FunctionName,
		Name:            aws.String(CurrentAlias),
		FunctionVersion: updated.Version,
	})

	if err != nil {
		return nil
	}

	f.Log.WithFields(log.Fields{
		"version": *updated.Version,
		"name":    f.FunctionName,
	}).Info("deployed")

	return nil
}

// Create the function with the given `zip`.
func (f *Function) Create(zip []byte) error {
	f.Log.Info("creating function")

	created, err := f.Service.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: &f.FunctionName,
		Description:  &f.Description,
		MemorySize:   &f.Memory,
		Timeout:      &f.Timeout,
		Runtime:      aws.String(f.runtime.Name()),
		Handler:      aws.String(f.runtime.Handler()),
		Role:         aws.String(f.Role),
		Publish:      aws.Bool(true),
		Code: &lambda.FunctionCode{
			ZipFile: zip,
		},
	})

	if err != nil {
		return err
	}

	f.Log.Info("creating alias")

	_, err = f.Service.CreateAlias(&lambda.CreateAliasInput{
		FunctionName:    &f.FunctionName,
		FunctionVersion: created.Version,
		Name:            aws.String(CurrentAlias),
	})

	if err != nil {
		return nil
	}

	f.Log.WithFields(log.Fields{
		"version": *created.Version,
		"name":    f.FunctionName,
	}).Info("deployed")

	return nil
}

// Invoke the remote Lambda function, returning the response and logs, if any.
func (f *Function) Invoke(event, context interface{}, kind InvocationType) (reply, logs io.Reader, err error) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, nil, err
	}

	contextBytes, err := json.Marshal(context)
	if err != nil {
		return nil, nil, err
	}

	res, err := f.Service.Invoke(&lambda.InvokeInput{
		ClientContext:  aws.String(base64.StdEncoding.EncodeToString(contextBytes)),
		FunctionName:   &f.FunctionName,
		InvocationType: aws.String(string(kind)),
		LogType:        aws.String("Tail"),
		Qualifier:      aws.String(CurrentAlias),
		Payload:        eventBytes,
	})

	if err != nil {
		return nil, nil, err
	}

	if res.FunctionError != nil {
		e := &InvokeError{
			Handled: *res.FunctionError == "Handled",
		}

		if err := json.Unmarshal(res.Payload, e); err != nil {
			return nil, nil, err
		}

		return nil, nil, e
	}

	if kind == Event {
		return bytes.NewReader(nil), bytes.NewReader(nil), nil
	}

	logs = base64.NewDecoder(base64.StdEncoding, strings.NewReader(*res.LogResult))
	reply = bytes.NewReader(res.Payload)
	return reply, logs, nil
}

// Rollback the function to the previous.
func (f *Function) Rollback() error {
	f.Log.Info("rolling back")

	alias, err := f.Service.GetAlias(&lambda.GetAliasInput{
		FunctionName: &f.FunctionName,
		Name:         aws.String(CurrentAlias),
	})

	if err != nil {
		return err
	}

	f.Log.Infof("current version: %s", *alias.FunctionVersion)

	list, err := f.Service.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: &f.FunctionName,
	})

	if err != nil {
		return err
	}

	versions := list.Versions[1:] // remove $LATEST
	if len(versions) < 2 {
		return errors.New("Can't rollback. Only one version deployed.")
	}

	latest := *versions[len(versions)-1].Version
	prev := *versions[len(versions)-2].Version
	rollback := latest

	if *alias.FunctionVersion == latest {
		rollback = prev
	}

	f.Log.Infof("rollback to version: %s", rollback)

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &f.FunctionName,
		Name:            aws.String(CurrentAlias),
		FunctionVersion: &rollback,
	})

	return err
}

// RollbackVersion the function to the specified version.
func (f *Function) RollbackVersion(version string) error {
	f.Log.Info("rolling back")

	alias, err := f.Service.GetAlias(&lambda.GetAliasInput{
		FunctionName: &f.FunctionName,
		Name:         aws.String(CurrentAlias),
	})

	if err != nil {
		return err
	}

	f.Log.Infof("current version: %s", *alias.FunctionVersion)

	if version == *alias.FunctionVersion {
		return errors.New("Specified version currently deployed.")
	}

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &f.FunctionName,
		Name:            aws.String(CurrentAlias),
		FunctionVersion: &version,
	})

	return err
}

// BuildBytes returns the generated zip as bytes.
func (f *Function) BuildBytes() ([]byte, error) {
	r, err := f.Build()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f.Log.Infof("created build (%s)", humanize.Bytes(uint64(len(b))))
	return b, nil
}

// Build returns the zipped contents of the function.
func (f *Function) Build() (io.Reader, error) {
	f.Log.Debugf("creating build")

	if err := f.RunHook("build"); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zip := archive.NewZipWriter(buf)

	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		f.Log.Debugf("compiling")
		if err := r.Build(f.Path); err != nil {
			return nil, fmt.Errorf("compiling: %s", err)
		}
	}

	if f.env != nil {
		f.Log.Debugf("adding .env.json")

		b, err := json.Marshal(f.env)
		if err != nil {
			return nil, err
		}

		zip.AddBytes(".env.json", b)
	}

	if f.runtime.Shimmed() {
		f.Log.Debugf("adding nodejs shim")
		zip.AddBytes("index.js", shim.MustAsset("index.js"))
		zip.AddBytes("byline.js", shim.MustAsset("byline.js"))
	}

	files, err := utils.LoadFiles(f.Path, f.IgnoredPatterns)
	if err != nil {
		return nil, err
	}

	for path, file := range files {
		f.Log.WithField("file", path).Debug("add file")
		if err := zip.AddFile(path, file); err != nil {
			return nil, err
		}
		defer file.Close()
	}

	if err := zip.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// Clean removes build artifacts from compiled runtimes.
func (f *Function) Clean() error {
	if err := f.RunHook("clean"); err != nil {
		return err
	}

	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		return r.Clean(f.Path)
	}

	return nil
}

// A HookError represents a failed hook command.
type HookError struct {
	Name    string
	Command string
	Output  string
}

// Error string.
func (e *HookError) Error() string {
	return fmt.Sprintf("hook %q: %s", e.Name, e.Output)
}

// RunHook executes hook `name` in the function's directory.
func (f *Function) RunHook(name string) error {
	var command string

	switch name {
	case "clean":
		command = f.Hooks.Clean
	case "build":
		command = f.Hooks.Build
	}

	if command == "" {
		return nil
	}

	f.Log.WithFields(log.Fields{
		"hook":    name,
		"command": command,
	}).Debug("hook")

	cmd := exec.Command("sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("FUNCTION=%s", f.Name))
	cmd.Dir = f.Path

	b, err := cmd.CombinedOutput()
	if err != nil {
		return &HookError{
			Name:    name,
			Command: command,
			Output:  string(b),
		}
	}

	return nil
}
