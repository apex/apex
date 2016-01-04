// Package function implements higher-level functionality for dealing with Lambda functions.
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

	"github.com/apex/apex/runtime"
	"github.com/apex/apex/shim"
	"github.com/apex/apex/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/dustin/go-humanize"
	"github.com/jpillora/archive"
)

// Errors.
var (
	ErrUnchanged = errors.New("function: unchanged")
)

// InvocationType determines how an invocation request is made.
type InvocationType string

// Invocation types.
const (
	RequestResponse InvocationType = "RequestResponse"
	Event                          = "Event"
	DryRun                         = "DryRun"
)

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
	Name        string `json:"name"`
	Description string `json:"description"`
	Runtime     string `json:"runtime"`
	Memory      int64  `json:"memory"`
	Timeout     int64  `json:"timeout"`
	Role        string `json:"role"`
	Main        string `json:"main"`
}

// Function represents a Lambda function, with configuration loaded
// from the "lambda.json" file on disk. Operations are performed
// against the function directory as the CWD, so os.Chdir() first.
type Function struct {
	Config
	Path    string
	Verbose bool
	Service lambdaiface.LambdaAPI
	runtime runtime.Runtime
	env     map[string]string
}

// log output when verbose output is enabled.
func (f *Function) log(msg string) {
	if f.Verbose {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", f.Name, msg)
	}
}

// logf formats message and output when verbose output is enabled.
func (f *Function) logf(msg string, v ...interface{}) {
	if f.Verbose {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", f.Name, fmt.Sprintf(msg, v...))
	}
}

// Open the lambda.json file and prime the config.
func (f *Function) Open() error {
	p, err := os.Open(filepath.Join(f.Path, "lambda.json"))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(p).Decode(&f.Config); err != nil {
		return err
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

// Deploy generates a zip and creates or updates the function.
func (f *Function) Deploy() error {
	f.log("deploying")

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
		f.log("unchanged")
		return ErrUnchanged
	}

	return f.Update(zip)
}

// DeployConfig deploys changes to configuration.
func (f *Function) DeployConfig() error {
	f.log("deploying config")
	_, err := f.Service.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: &f.Name,
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
	f.log("deleting")

	_, err := f.Service.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: &f.Name,
	})

	return err
}

// Info returns the function information.
func (f *Function) Info() (*lambda.GetFunctionOutput, error) {
	f.log("fetching config")
	return f.Service.GetFunction(&lambda.GetFunctionInput{
		FunctionName: &f.Name,
	})
}

// Update the function with the given `zip`.
func (f *Function) Update(zip []byte) error {
	f.log("updating")

	_, err := f.Service.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: &f.Name,
		Publish:      aws.Bool(true),
		ZipFile:      zip,
	})

	return err
}

// Create the function with the given `zip`.
func (f *Function) Create(zip []byte) error {
	f.log("creating")

	_, err := f.Service.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: &f.Name,
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
		FunctionName:   aws.String(f.Name),
		InvocationType: aws.String(string(kind)),
		LogType:        aws.String("Tail"),
		Qualifier:      aws.String("$LATEST"),
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

// Clean removes build artifacts from compiled runtimes.
func (f *Function) Clean() error {
	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		return r.Clean()
	}
	return nil
}

// Zip returns the zipped contents of the function.
func (f *Function) Zip() (io.Reader, error) {
	buf := new(bytes.Buffer)
	zip := archive.NewZipWriter(buf)

	if r, ok := f.runtime.(runtime.CompiledRuntime); ok {
		if err := r.Build(f.Main); err != nil {
			return nil, fmt.Errorf("compiling: %s", err)
		}
	}

	if f.env != nil {
		b, err := json.Marshal(f.env)
		if err != nil {
			return nil, err
		}

		zip.AddBytes(".env.json", b)
	}

	if f.runtime.Shimmed() {
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
	f.log("creating zip")

	r, err := f.Zip()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f.logf("created zip (%s)", humanize.Bytes(uint64(len(b))))
	return b, nil
}
