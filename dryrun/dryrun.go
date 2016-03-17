// Package dryrun implements the Lambda API in order to no-op changes, and display dry-run output.
package dryrun

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dustin/go-humanize"

	"github.com/apex/apex/utils"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
)

// Lambda is a partially implemented Lambda API implementation used to perform a dry-run.
type Lambda struct {
	*lambda.Lambda
}

// New dry-run Lambda service for the given session.
func New(session *session.Session) *Lambda {
	fmt.Printf("\n")
	return &Lambda{
		Lambda: lambda.New(session),
	}
}

// CreateFunction stub.
func (l *Lambda) CreateFunction(in *lambda.CreateFunctionInput) (*lambda.FunctionConfiguration, error) {
	l.create("function", *in.FunctionName, map[string]interface{}{
		"runtime": *in.Runtime,
		"memory":  *in.MemorySize,
		"timeout": *in.Timeout,
		"handler": *in.Handler,
	})

	out := &lambda.FunctionConfiguration{
		Version: aws.String("1"),
	}

	return out, nil
}

// UpdateFunctionCode stub.
func (l *Lambda) UpdateFunctionCode(in *lambda.UpdateFunctionCodeInput) (*lambda.FunctionConfiguration, error) {
	res, err := l.GetFunction(&lambda.GetFunctionInput{
		FunctionName: in.FunctionName,
	})

	if err != nil {
		return nil, err
	}

	size := uint64(len(in.ZipFile))
	checksum := utils.Sha256(in.ZipFile)
	remoteChecksum := *res.Configuration.CodeSha256
	remoteSize := uint64(*res.Configuration.CodeSize)

	if checksum != remoteChecksum {
		l.create("function", *in.FunctionName, map[string]interface{}{
			"size": fmt.Sprintf("%s -> %s", humanize.Bytes(remoteSize), humanize.Bytes(size)),
		})
	}

	out := &lambda.FunctionConfiguration{
		Version: aws.String("$LATEST"),
	}

	return out, nil
}

// UpdateFunctionConfiguration stub.
func (l *Lambda) UpdateFunctionConfiguration(in *lambda.UpdateFunctionConfigurationInput) (*lambda.FunctionConfiguration, error) {
	res, err := l.GetFunctionConfiguration(&lambda.GetFunctionConfigurationInput{
		FunctionName: in.FunctionName,
	})

	if e, ok := err.(awserr.Error); ok {
		if e.Code() == "ResourceNotFoundException" {
			return nil, nil
		}
	}

	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})

	if *in.Description != *res.Description {
		m["description"] = fmt.Sprintf("%q -> %q", *res.Description, *in.Description)
	}

	if *in.Handler != *res.Handler {
		m["handler"] = fmt.Sprintf("%s -> %s", *res.Handler, *in.Handler)
	}

	if *in.MemorySize != *res.MemorySize {
		m["memory"] = fmt.Sprintf("%v -> %v", *res.MemorySize, *in.MemorySize)
	}

	if *in.Role != *res.Role {
		m["role"] = fmt.Sprintf("%v -> %v", *res.Role, *in.Role)
	}

	if *in.Timeout != *res.Timeout {
		m["timeout"] = fmt.Sprintf("%v -> %v", *res.Timeout, *in.Timeout)
	}

	if len(m) > 0 {
		l.update("config", *in.FunctionName, m)
	}

	return nil, nil
}

// DeleteFunction stub.
func (l *Lambda) DeleteFunction(in *lambda.DeleteFunctionInput) (*lambda.DeleteFunctionOutput, error) {
	if in.Qualifier == nil {
		l.remove("function", *in.FunctionName, nil)
	} else {
		l.remove("function version", fmt.Sprintf("%s (version: %s)", *in.FunctionName, *in.Qualifier), nil)
	}

	return nil, nil
}

// CreateAlias stub.
func (l *Lambda) CreateAlias(in *lambda.CreateAliasInput) (*lambda.AliasConfiguration, error) {
	l.create("alias", *in.FunctionName, map[string]interface{}{
		"alias":   *in.Name,
		"version": *in.FunctionVersion,
	})
	return nil, nil
}

// UpdateAlias stub.
func (l *Lambda) UpdateAlias(in *lambda.UpdateAliasInput) (*lambda.AliasConfiguration, error) {
	l.update("alias", *in.FunctionName, map[string]interface{}{
		"alias":   *in.Name,
		"version": *in.FunctionVersion,
	})
	return nil, nil
}

func (l *Lambda) log(kind, name string, m map[string]interface{}, symbol rune, color int) {
	fmt.Printf("  \033[%dm%c %s\033[0m \033[%dm%s\033[0m\n", color, symbol, kind, blue, name)
	for k, v := range m {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Printf("\n")
}

// create message.
func (l *Lambda) create(kind, name string, m map[string]interface{}) {
	l.log(kind, name, m, '+', green)
}

// update message.
func (l *Lambda) update(kind, name string, m map[string]interface{}) {
	l.log(kind, name, m, '~', yellow)
}

// remove message.
func (l *Lambda) remove(kind, name string, m map[string]interface{}) {
	l.log(kind, name, m, '-', red)
}
