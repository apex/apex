// Package dryrun implements the Lambda API in order to no-op changes, and display dry-run output.
package dryrun

import (
	"fmt"

	"github.com/apex/apex/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dustin/go-humanize"
)

const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

type change struct {
	Name string
	From interface{}
	To   interface{}
}

// TODO(tj): sync so concurrent writes don't race
// TODO(tj): spend more time on nicer output

// Lambda is a partially implemented Lambda API implementation used to perform a dry-run.
type Lambda struct {
	*lambda.Lambda
}

// New dry-run Lambda service for the given session.
func New(session *session.Session) *Lambda {
	return &Lambda{
		Lambda: lambda.New(session),
	}
}

// CreateFunction stub.
func (l *Lambda) CreateFunction(in *lambda.CreateFunctionInput) (*lambda.FunctionConfiguration, error) {
	l.createf("function", *in.FunctionName, "runtime=%s memory=%d timeout=%d handler=%s", *in.Runtime, *in.MemorySize, *in.Timeout, *in.Handler)

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
		l.updatef("function", *in.FunctionName, "%s -> %s", humanize.Bytes(size), humanize.Bytes(remoteSize))
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

	var changes []change

	if *in.Description != *res.Description {
		changes = append(changes, change{
			Name: "description",
			From: *res.Description,
			To:   *in.Description,
		})
	}

	if *in.Handler != *res.Handler {
		changes = append(changes, change{
			Name: "handler",
			From: *res.Handler,
			To:   *in.Handler,
		})
	}

	if *in.MemorySize != *res.MemorySize {
		changes = append(changes, change{
			Name: "memory",
			From: *res.MemorySize,
			To:   *in.MemorySize,
		})
	}

	if *in.Role != *res.Role {
		changes = append(changes, change{
			Name: "role",
			From: *res.Role,
			To:   *in.Role,
		})
	}

	if *in.Timeout != *res.Timeout {
		changes = append(changes, change{
			Name: "timeout",
			From: *res.Timeout,
			To:   *in.Timeout,
		})
	}

	if len(changes) > 0 {
		l.updatef("config", *in.FunctionName, "")
		for _, change := range changes {
			switch change.From.(type) {
			case string:
				l.updatef("config", change.Name, "%q -> %q", change.From, change.To)
			default:
				l.updatef("config", change.Name, "%v -> %v", change.From, change.To)
			}
		}
	}

	return nil, nil
}

// DeleteFunction stub.
func (l *Lambda) DeleteFunction(in *lambda.DeleteFunctionInput) (*lambda.DeleteFunctionOutput, error) {
	l.removef("function", *in.FunctionName, "")
	return nil, nil
}

// CreateAlias stub.
func (l *Lambda) CreateAlias(in *lambda.CreateAliasInput) (*lambda.AliasConfiguration, error) {
	l.createf("alias", *in.FunctionName, "%s -> %s", *in.Name, *in.FunctionVersion)
	return nil, nil
}

// UpdateAlias stub.
func (l *Lambda) UpdateAlias(in *lambda.UpdateAliasInput) (*lambda.AliasConfiguration, error) {
	l.updatef("alias", *in.FunctionName, "%s -> %s", *in.Name, *in.FunctionVersion)
	return nil, nil
}

// create message.
func (l *Lambda) createf(key, name, msg string, args ...interface{}) {
	fmt.Printf("  \033[%dm+ %-10s\033[0m \033[%dm%-13s\033[0m %s\n", green, key, blue, name, fmt.Sprintf(msg, args...))
}

// update message.
func (l *Lambda) updatef(key, name, msg string, args ...interface{}) {
	fmt.Printf("  \033[%dm~ %-10s\033[0m \033[%dm%-13s\033[0m %s\n", yellow, key, blue, name, fmt.Sprintf(msg, args...))
}

// remove message.
func (l *Lambda) removef(key, name, msg string, args ...interface{}) {
	fmt.Printf("  \033[%dm- %-10s\033[0m \033[%dm%-13s\033[0m %s\n", red, key, blue, name, fmt.Sprintf(msg, args...))
}
