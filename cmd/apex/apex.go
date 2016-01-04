package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/apex/apex/function"
	"github.com/apex/apex/logs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/segmentio/go-prompt"
	"github.com/tj/docopt"
)

var version = "0.1.0"

const usage = `
  Usage:
    apex deploy [-C path] [--env name=val]... [-v]
    apex delete [-C path] [-y] [-v]
    apex invoke [-C path] [--async] [-v]
    apex build [-C path] [-v]
    apex logs
    apex -h | --help
    apex --version

  Options:
    -e, --env name=val  Environment variable
    -a, --async         Async invocation
    -C, --chdir path    Working directory
    -y, --yes           Automatic yes to prompts
    -h, --help          Output help information
    -v, --verbose       Output verbose logs
    -V, --version       Output version

  Examples:
    Deploy a function in the current directory
    $ apex deploy

    Delete a function in the current directory
    $ apex delete

    Invoke a function in the current directory
    $ apex invoke < request.json

    Deploy a function in a different directory
    $ apex deploy -C functions/hello-world

    Output zip of a function in the current directory
    $ apex build > /tmp/out.zip

    Tail Kinesis logs
    $ apex logs
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	switch {
	case args["logs"].(bool):
		tailLogs()
		return
	}

	fn := &function.Function{
		Service: lambda.New(session.New(aws.NewConfig())),
		Verbose: args["--verbose"].(bool),
		Path:    ".",
	}

	if dir, ok := args["--chdir"].(string); ok {
		if err := os.Chdir(dir); err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	if err := fn.Open(); err != nil {
		log.Fatalf("error: %s", err)
	}

	switch {
	case args["deploy"].(bool):
		deploy(fn, args["--env"].([]string))
	case args["delete"].(bool):
		if args["--yes"].(bool) || prompt.Confirm("Are you sure? [yes/no]") {
			delete(fn)
		}
	case args["invoke"].(bool):
		invoke(fn, args["--verbose"].(bool), args["--async"].(bool))
	case args["build"].(bool):
		build(fn)
	}
}

// invoke reads request json from stdin and outputs the responses.
func invoke(fn *function.Function, verbose, async bool) {
	dec := json.NewDecoder(os.Stdin)
	kind := function.RequestResponse

	if async {
		kind = function.Event
	}

	for {
		var v struct {
			Event   interface{}
			Context interface{}
		}

		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error: %s", err)
		}

		reply, logs, err := fn.Invoke(v.Event, v.Context, kind)
		if err != nil {
			log.Fatalf("error: %s", err)
		}

		if verbose {
			io.Copy(os.Stderr, logs)
		}

		io.Copy(os.Stdout, reply)
		fmt.Fprintf(os.Stdout, "\n")
	}
}

// deploy creates or updates the function.
func deploy(fn *function.Function, env []string) {
	for _, s := range env {
		parts := strings.Split(s, "=")
		fn.SetEnv(parts[0], parts[1])
	}

	if err := fn.Deploy(); err != nil && err != function.ErrUnchanged {
		log.Fatalf("error deploying code: %s", err)
	}

	if err := fn.DeployConfig(); err != nil {
		log.Fatalln("error deploying config: %s", err)
	}

	fn.Clean()
}

// delete the function.
func delete(fn *function.Function) {
	if err := fn.Delete(); err != nil {
		log.Fatalf("error: %s", err)
	}
}

// build outputs the generated archive to stdout.
func build(fn *function.Function) {
	zip, err := fn.Zip()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	_, err = io.Copy(os.Stdout, zip)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// tail Kinesis logs for changes.
// TODO(tj): parse json for display via apex/log
func tailLogs() {
	client := kinesis.New(session.New(aws.NewConfig()))

	tailer := logs.Tailer{
		Stream:       "logs",
		Service:      client,
		PollInterval: 500 * time.Millisecond,
	}

	for record := range tailer.Start() {
		fmt.Printf("%s\n", record.Data)
	}

	if err := tailer.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
