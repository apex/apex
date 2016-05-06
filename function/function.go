// Package function implements function-level operations.
package function

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/dustin/go-humanize"
	"github.com/jpillora/archive"
	"gopkg.in/validator.v2"

	"github.com/apex/apex/hooks"
	"github.com/apex/apex/utils"
	"github.com/apex/apex/vpc"
)

// timelessInfo is used to zero mtime which causes function checksums
// to change regardless of the contents actually being altered, specifically
// when using tools such as browserify or webpack.
type timelessInfo struct {
	os.FileInfo
}

func (t timelessInfo) ModTime() time.Time {
	return time.Unix(0, 0)
}

// defaultPlugins are the default plugins which are required by Apex. Note that
// the order here is important for some plugins such as inference before the
// runtimes.
var defaultPlugins = []string{
	"inference",
	"golang",
	"python",
	"nodejs",
	"java",
	"hooks",
	"env",
	"shim",
}

// InvocationType determines how an invocation request is made.
type InvocationType string

// Invocation types.
const (
	RequestResponse InvocationType = "RequestResponse"
	Event                          = "Event"
	DryRun                         = "DryRun"
)

// CurrentAlias name.
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
	Description      string            `json:"description"`
	Runtime          string            `json:"runtime" validate:"nonzero"`
	Memory           int64             `json:"memory" validate:"nonzero"`
	Timeout          int64             `json:"timeout" validate:"nonzero"`
	Role             string            `json:"role" validate:"nonzero"`
	Handler          string            `json:"handler" validate:"nonzero"`
	Shim             bool              `json:"shim"`
	Environment      map[string]string `json:"environment"`
	Hooks            hooks.Hooks       `json:"hooks"`
	RetainedVersions int               `json:"retainedVersions"`
	VPC              vpc.VPC           `json:"vpc"`
}

// Function represents a Lambda function, with configuration loaded
// from the "function.json" file on disk.
type Function struct {
	Config
	Name         string
	FunctionName string
	Path         string
	Service      lambdaiface.LambdaAPI
	Log          log.Interface
	IgnoreFile   []byte
	Plugins      []string
	Alias        string
}

// Open the function.json file and prime the config.
func (f *Function) Open() error {
	f.defaults()
	f.Log = f.Log.WithField("function", f.Name)

	p, err := os.Open(filepath.Join(f.Path, "function.json"))
	if err == nil {
		if err := json.NewDecoder(p).Decode(&f.Config); err != nil {
			return err
		}
	}

	if err := f.hookOpen(); err != nil {
		return err
	}

	if err := validator.Validate(&f.Config); err != nil {
		return fmt.Errorf("error opening function %s: %s", f.Name, err.Error())
	}

	ignoreFile, err := utils.ReadIgnoreFile(f.Path)
	if err != nil {
		return err
	}
	f.IgnoreFile = append(f.IgnoreFile, []byte("\n")...)
	f.IgnoreFile = append(f.IgnoreFile, ignoreFile...)

	return nil
}

// defaults applies configuration defaults.
func (f *Function) defaults() {
	if f.Alias == "" {
		f.Alias = CurrentAlias
	}

	if f.Plugins == nil {
		f.Plugins = defaultPlugins
	}

	if f.Environment == nil {
		f.Environment = make(map[string]string)
	}

	if f.VPC.Subnets == nil {
		f.VPC.Subnets = []string{}
	}

	if f.VPC.SecurityGroups == nil {
		f.VPC.SecurityGroups = []string{}
	}

	f.Setenv("APEX_FUNCTION_NAME", f.Name)
	f.Setenv("LAMBDA_FUNCTION_NAME", f.FunctionName)
}

// Setenv sets environment variable `name` to `value`.
func (f *Function) Setenv(name, value string) {
	f.Environment[name] = value
}

// Deploy generates a zip and creates or deploy the function.
// If the configuration hasn't been changed it will deploy only code,
// otherwise it will deploy both configuration and code.
func (f *Function) Deploy() error {
	f.Log.Debug("deploying")

	zip, err := f.BuildBytes()
	if err != nil {
		return err
	}

	if err := f.hookDeploy(); err != nil {
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

	if f.configChanged(config) {
		f.Log.Debug("config changed")
		return f.DeployConfigAndCode(zip)
	}

	f.Log.Info("config unchanged")
	return f.DeployCode(zip, config)
}

// DeployCode deploys function code when changed.
func (f *Function) DeployCode(zip []byte, config *lambda.GetFunctionOutput) error {
	remoteHash := *config.Configuration.CodeSha256
	localHash := utils.Sha256(zip)

	if localHash == remoteHash {
		f.Log.Info("code unchanged")
		return nil
	}

	f.Log.WithFields(log.Fields{
		"local":  localHash,
		"remote": remoteHash,
	}).Debug("code changed")

	return f.Update(zip)
}

// DeployConfigAndCode updates config and updates function code.
func (f *Function) DeployConfigAndCode(zip []byte) error {
	f.Log.Info("updating config")

	_, err := f.Service.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: &f.FunctionName,
		MemorySize:   &f.Memory,
		Timeout:      &f.Timeout,
		Description:  &f.Description,
		Role:         &f.Role,
		Runtime:      &f.Runtime,
		Handler:      &f.Handler,
		VpcConfig: &lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(f.VPC.SecurityGroups),
			SubnetIds:        aws.StringSlice(f.VPC.Subnets),
		},
	})

	if err != nil {
		return err
	}

	return f.Update(zip)
}

// Delete the function including all its versions
func (f *Function) Delete() error {
	f.Log.Info("deleting")
	_, err := f.Service.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: &f.FunctionName,
	})

	if err != nil {
		return err
	}

	f.Log.Info("function deleted")

	return nil
}

// GetConfig returns the function configuration.
func (f *Function) GetConfig() (*lambda.GetFunctionOutput, error) {
	f.Log.Debug("fetching config")
	return f.Service.GetFunction(&lambda.GetFunctionInput{
		FunctionName: &f.FunctionName,
	})
}

// GetConfigQualifier returns the function configuration for the given qualifier.
func (f *Function) GetConfigQualifier(s string) (*lambda.GetFunctionOutput, error) {
	f.Log.Debug("fetching config")
	return f.Service.GetFunction(&lambda.GetFunctionInput{
		FunctionName: &f.FunctionName,
		Qualifier:    &s,
	})
}

// GetConfigCurrent returns the function configuration for the current version.
func (f *Function) GetConfigCurrent() (*lambda.GetFunctionOutput, error) {
	return f.GetConfigQualifier(f.Alias)
}

// Update the function with the given `zip`.
func (f *Function) Update(zip []byte) error {
	f.Log.Info("updating function")

	versionsToCleanup, err := f.versionsToCleanup()
	if err != nil {
		return err
	}

	updated, err := f.Service.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: &f.FunctionName,
		Publish:      aws.Bool(true),
		ZipFile:      zip,
	})

	if err != nil {
		return err
	}

	if err := f.CreateOrUpdateAlias(f.Alias, *updated.Version); err != nil {
		return err
	}

	f.Log.WithFields(log.Fields{
		"version": *updated.Version,
		"name":    f.FunctionName,
	}).Info("function updated")

	return f.removeVersions(versionsToCleanup)
}

// Create the function with the given `zip`.
func (f *Function) Create(zip []byte) error {
	f.Log.Info("creating function")

	created, err := f.Service.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: &f.FunctionName,
		Description:  &f.Description,
		MemorySize:   &f.Memory,
		Timeout:      &f.Timeout,
		Runtime:      &f.Runtime,
		Handler:      &f.Handler,
		Role:         &f.Role,
		Publish:      aws.Bool(true),
		Code: &lambda.FunctionCode{
			ZipFile: zip,
		},
		VpcConfig: &lambda.VpcConfig{
			SecurityGroupIds: aws.StringSlice(f.VPC.SecurityGroups),
			SubnetIds:        aws.StringSlice(f.VPC.Subnets),
		},
	})

	if err != nil {
		return err
	}

	if err := f.CreateOrUpdateAlias(f.Alias, *created.Version); err != nil {
		return err
	}

	f.Log.WithFields(log.Fields{
		"version": *created.Version,
		"name":    f.FunctionName,
	}).Info("function created")

	return nil
}

// CreateOrUpdateAlias attempts creating the alias, or updates if it already exists.
func (f *Function) CreateOrUpdateAlias(alias, version string) error {
	_, err := f.Service.CreateAlias(&lambda.CreateAliasInput{
		FunctionName:    &f.FunctionName,
		FunctionVersion: &version,
		Name:            &alias,
	})

	if err == nil {
		f.Log.WithField("version", version).Infof("created alias %s", alias)
		return nil
	}

	if e, ok := err.(awserr.Error); !ok || e.Code() != "ResourceConflictException" {
		return err
	}

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &f.FunctionName,
		FunctionVersion: &version,
		Name:            &alias,
	})

	if err != nil {
		return err
	}

	f.Log.WithField("version", version).Infof("updated alias %s", alias)
	return nil
}

// Invoke the remote Lambda function, returning the response and logs, if any.
func (f *Function) Invoke(event, context interface{}) (reply, logs io.Reader, err error) {
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
		InvocationType: aws.String(string(RequestResponse)),
		LogType:        aws.String("Tail"),
		Qualifier:      &f.Alias,
		Payload:        eventBytes,
	})

	if err != nil {
		return nil, nil, err
	}

	logs = base64.NewDecoder(base64.StdEncoding, strings.NewReader(*res.LogResult))

	if res.FunctionError != nil {
		e := &InvokeError{
			Handled: *res.FunctionError == "Handled",
		}

		if err := json.Unmarshal(res.Payload, e); err != nil {
			return nil, logs, err
		}

		return nil, logs, e
	}

	reply = bytes.NewReader(res.Payload)
	return reply, logs, nil
}

// Rollback the function to the previous.
func (f *Function) Rollback() error {
	f.Log.Info("rolling back")

	alias, err := f.currentVersionAlias()
	if err != nil {
		return err
	}

	f.Log.Debugf("current version: %s", *alias.FunctionVersion)

	versions, err := f.versions()
	if err != nil {
		return err
	}

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
		Name:            &f.Alias,
		FunctionVersion: &rollback,
	})

	if err != nil {
		return err
	}

	f.Log.WithField("current version", rollback).Info("function rolled back")

	return nil
}

// RollbackVersion the function to the specified version.
func (f *Function) RollbackVersion(version string) error {
	f.Log.Info("rolling back")

	alias, err := f.currentVersionAlias()
	if err != nil {
		return err
	}

	f.Log.Debugf("current version: %s", *alias.FunctionVersion)

	if version == *alias.FunctionVersion {
		return errors.New("Specified version currently deployed.")
	}

	f.Log.Infof("rollback to version: %s", version)

	_, err = f.Service.UpdateAlias(&lambda.UpdateAliasInput{
		FunctionName:    &f.FunctionName,
		Name:            &f.Alias,
		FunctionVersion: &version,
	})

	if err != nil {
		return err
	}

	f.Log.WithField("current version", version).Info("function rolled back")

	return nil
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

	f.Log.Debugf("created build (%s)", humanize.Bytes(uint64(len(b))))
	return b, nil
}

// Build returns the zipped contents of the function.
func (f *Function) Build() (io.Reader, error) {
	f.Log.Debugf("creating build")

	buf := new(bytes.Buffer)
	zip := archive.NewCompressedZipWriter(buf, flate.DefaultCompression)

	if err := f.hookBuild(zip); err != nil {
		return nil, err
	}

	files, err := utils.LoadFiles(f.Path, f.IgnoreFile)
	if err != nil {
		return nil, err
	}

	for _, path := range files {
		f.Log.WithField("file", path).Debug("add file to zip")

		info, err := os.Lstat(filepath.Join(f.Path, path))
		if err != nil {
			f.Log.Debugf("Lstat error: %v", err)
			return nil, err
		}

		realPath := filepath.Join(f.Path, path)
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			linkPath, err := filepath.EvalSymlinks(filepath.Join(f.Path, path))
			if err != nil {
				f.Log.Debugf("EvalSymlinks error: %v", err)
				return nil, err
			}
			realPath = linkPath
		}

		fh, err := os.Open(realPath)
		if err != nil {
			f.Log.Debugf("Open error: %v", err)
			return nil, err
		}

		info, err = fh.Stat()
		if err != nil {
			return nil, err
		}

		unixPath := strings.Replace(path, "\\", "/", -1)
		if err := zip.AddInfoFile(unixPath, timelessInfo{info}, fh); err != nil {
			return nil, err
		}

		if err := fh.Close(); err != nil {
			return nil, err
		}
	}

	if err := zip.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// Clean invokes the CleanHook, useful for removing build artifacts and so on.
func (f *Function) Clean() error {
	return f.hookClean()
}

// GroupName returns the CloudWatchLogs group name.
func (f *Function) GroupName() string {
	return fmt.Sprintf("/aws/lambda/%s", f.FunctionName)
}

// versions returns list of all versions deployed to AWS Lambda
func (f *Function) versions() ([]*lambda.FunctionConfiguration, error) {
	list, err := f.Service.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: &f.FunctionName,
	})

	if err != nil {
		return nil, err
	}

	versions := list.Versions[1:] // remove $LATEST

	return versions, nil
}

// versionsToCleanup returns list of versions to remove after updating function
func (f *Function) versionsToCleanup() ([]*lambda.FunctionConfiguration, error) {
	versions, err := f.versions()
	if err != nil {
		return nil, err
	}

	if len(versions) > f.RetainedVersions {
		return versions[:len(versions)-f.RetainedVersions], nil
	}

	return nil, nil
}

// removeVersions removes specifed function's versions
func (f *Function) removeVersions(versions []*lambda.FunctionConfiguration) error {
	for _, v := range versions {
		f.Log.Debugf("cleaning up version: %s", *v.Version)

		_, err := f.Service.DeleteFunction(&lambda.DeleteFunctionInput{
			FunctionName: &f.FunctionName,
			Qualifier:    v.Version,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// currentVersionAlias returns alias configuration for currently deployed function
func (f *Function) currentVersionAlias() (*lambda.AliasConfiguration, error) {
	return f.Service.GetAlias(&lambda.GetAliasInput{
		FunctionName: &f.FunctionName,
		Name:         &f.Alias,
	})
}

// configChanged checks if function configuration differs from configuration stored in AWS Lambda
func (f *Function) configChanged(config *lambda.GetFunctionOutput) bool {
	type diffConfig struct {
		Description string
		Memory      int64
		Timeout     int64
		Role        string
		Runtime     string
		Handler     string
		VPC         vpc.VPC
	}

	localConfig := &diffConfig{
		Description: f.Description,
		Memory:      f.Memory,
		Timeout:     f.Timeout,
		Role:        f.Role,
		Runtime:     f.Runtime,
		Handler:     f.Handler,
		VPC: vpc.VPC{
			Subnets:        f.VPC.Subnets,
			SecurityGroups: f.VPC.SecurityGroups,
		},
	}

	remoteConfig := &diffConfig{
		Description: *config.Configuration.Description,
		Memory:      *config.Configuration.MemorySize,
		Timeout:     *config.Configuration.Timeout,
		Role:        *config.Configuration.Role,
		Runtime:     *config.Configuration.Runtime,
		Handler:     *config.Configuration.Handler,
	}

	// SDK is inconsistent here. VpcConfig can be nil or empty struct.
	remoteConfig.VPC = vpc.VPC{Subnets: []string{}, SecurityGroups: []string{}}
	if config.Configuration.VpcConfig != nil {
		remoteConfig.VPC = vpc.VPC{
			Subnets:        aws.StringValueSlice(config.Configuration.VpcConfig.SubnetIds),
			SecurityGroups: aws.StringValueSlice(config.Configuration.VpcConfig.SecurityGroupIds),
		}
	}

	localConfigJSON, _ := json.Marshal(localConfig)
	remoteConfigJSON, _ := json.Marshal(remoteConfig)
	return string(localConfigJSON) != string(remoteConfigJSON)
}

// hookOpen calls Openers.
func (f *Function) hookOpen() error {
	for _, name := range f.Plugins {
		if p, ok := plugins[name].(Opener); ok {
			if err := p.Open(f); err != nil {
				return err
			}
		}
	}
	return nil
}

// hookBuild calls Builders.
func (f *Function) hookBuild(zip *archive.Archive) error {
	for _, name := range f.Plugins {
		if p, ok := plugins[name].(Builder); ok {
			if err := p.Build(f, zip); err != nil {
				return err
			}
		}
	}
	return nil
}

// hookClean calls Cleaners.
func (f *Function) hookClean() error {
	for _, name := range f.Plugins {
		if p, ok := plugins[name].(Cleaner); ok {
			if err := p.Clean(f); err != nil {
				return err
			}
		}
	}
	return nil
}

// hookDeploy calls Deployers.
func (f *Function) hookDeploy() error {
	for _, name := range f.Plugins {
		if p, ok := plugins[name].(Deployer); ok {
			if err := p.Deploy(f); err != nil {
				return err
			}
		}
	}
	return nil
}
