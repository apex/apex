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
	"path/filepath"
	"strings"

	"gopkg.in/validator.v2"

	"github.com/apex/apex/runtime"
	"github.com/apex/apex/shim"
	"github.com/apex/apex/utils"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/dustin/go-humanize"
	"github.com/jpillora/archive"
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
	Name        string `json:"name" validate:"nonzero"`
	Description string `json:"description"`
	Runtime     string `json:"runtime" validate:"nonzero"`
	Memory      int64  `json:"memory" validate:"nonzero"`
	Timeout     int64  `json:"timeout" validate:"nonzero"`
	Role        string `json:"role" validate:"nonzero"`
}

// Function represents a Lambda function, with configuration loaded
// from the "function.json" file on disk. Operations are performed
// against the function directory as the CWD, so os.Chdir() first.
type Function struct {
	Config
	Path    string
	Prefix  string
	Service lambdaiface.LambdaAPI
	Log     log.Interface
	runtime runtime.Runtime
	env     map[string]string
}

// Open the function.json file and prime the config.
func (f *Function) Open() error {
	p, err := os.Open(filepath.Join(f.Path, "function.json"))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(p).Decode(&f.Config); err != nil {
		return err
	}

	if err := validator.Validate(&f.Config); err != nil {
		return fmt.Errorf("error opening function %s: %s", f.Name, err.Error())
	}

	r, err := runtime.ByName(f.Runtime)
	if err != nil {
		return err
	}
	f.runtime = r

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

	zip, err := f.ZipBytes()
	if err != nil {
		return err
	}

	info, err := f.Info()

	if e, ok := err.(awserr.Error); ok {
		if e.Code() == "ResourceNotFoundException" {
			return f.Create(zip)
		}
	}

	if err != nil {
		return err
	}

	remoteHash := *info.Configuration.CodeSha256
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
		FunctionName: f.name(),
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
		FunctionName: f.name(),
	})
	return err
}

// Info returns the function information.
func (f *Function) Info() (*lambda.GetFunctionOutput, error) {
	f.Log.Info("fetching config")
	return f.Service.GetFunction(&lambda.GetFunctionInput{
		FunctionName: f.name(),
	})
}

// Update the function with the given `zip`.
func (f *Function) Update(zip []byte) error {
	f.Log.Info("updating function")

	updated, err := f.Service.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: f.name(),
		Publish:      aws.Bool(true),
		ZipFile:      zip,
	})

	if err != nil {
		return err
	}

	f.Log.Info("updating alias")

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    f.name(),
		Name:            aws.String(CurrentAlias),
		FunctionVersion: updated.Version,
	})

	return err
}

// Create the function with the given `zip`.
func (f *Function) Create(zip []byte) error {
	f.Log.Info("creating function")

	created, err := f.Service.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: f.name(),
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
		FunctionName:    f.name(),
		FunctionVersion: created.Version,
		Name:            aws.String(CurrentAlias),
	})

	return err
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
		FunctionName:   f.name(),
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
		FunctionName: f.name(),
		Name:         aws.String(CurrentAlias),
	})
	if err != nil {
		return err
	}

	f.Log.Infof("current version: %s", *alias.FunctionVersion)

	list, err := f.Service.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: f.name(),
	})
	if err != nil {
		return err
	}

	versions := list.Versions
	_, versions = versions[0], versions[1:] // remove $LATEST
	if len(versions) < 2 {
		return errors.New("Can't rollback. Only one version deployed.")
	}
	latestVersion := *versions[len(versions)-1].Version
	prevVersion := *versions[len(versions)-2].Version

	rollbackToVersion := latestVersion
	if *alias.FunctionVersion == latestVersion {
		rollbackToVersion = prevVersion
	}

	f.Log.Infof("rollback to version: %s", rollbackToVersion)

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    f.name(),
		Name:            aws.String(CurrentAlias),
		FunctionVersion: &rollbackToVersion,
	})

	return err
}

// RollbackVersion the function to the specified version.
func (f *Function) RollbackVersion(version string) error {
	f.Log.Info("rolling back")

	alias, err := f.Service.GetAlias(&lambda.GetAliasInput{
		FunctionName: f.name(),
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
		FunctionName:    f.name(),
		Name:            aws.String(CurrentAlias),
		FunctionVersion: &version,
	})

	return err
}

// Clean removes build artifacts from compiled runtimes.
func (f *Function) Clean() error {
	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		return r.Clean(f.Path)
	}
	return nil
}

// Zip returns the zipped contents of the function.
func (f *Function) Zip() (io.Reader, error) {
	buf := new(bytes.Buffer)
	zip := archive.NewZipWriter(buf)

	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		f.Log.Debugf("compiling")
		if err := r.Build(f.Path); err != nil {
			return nil, fmt.Errorf("compiling: %s", err)
		}
	}

	// TODO(tj): remove or add --env flag back
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

	if err := zip.AddDir(f.Path); err != nil {
		return nil, err
	}

	if err := zip.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// ZipBytes returns the generated zip as bytes.
func (f *Function) ZipBytes() ([]byte, error) {
	f.Log.Debugf("creating zip")

	r, err := f.Zip()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f.Log.Infof("created zip (%s)", humanize.Bytes(uint64(len(b))))
	return b, nil
}

// name returns the canonical Lambda function name with optional prefix.
// This is used to prevent collisions with multiple projects using the
// same function names.
func (f *Function) name() *string {
	if f.Prefix == "" {
		return aws.String(f.Name)
	} else {
		return aws.String(fmt.Sprintf("%s_%s", f.Prefix, f.Name))
	}
}
