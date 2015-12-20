package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/apex/apex/function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/tj/docopt"
)

var version = "0.0.2"

const usage = `
  Usage:
    apex deploy [-C path] [--env name=val]...
    apex invoke [-C path] [--async] [-v]
    apex zip [-C path]
    apex -h | --help
    apex --version

  Options:
    -e, --env name=val  Environment variable
    -a, --async         Async invocation
    -C, --chdir path    Working directory
    -h, --help          Output help information
    -v, --verbose       Output verbose logs
    -V, --version       Output version

  Examples:
    Deploy a function in the current directory
    $ apex deploy

    Invoke a function in the current directory
    $ apex invoke < request.json

    Deploy a function in a different directory
    $ apex deploy -C functions/hello-world

    Output zip of a function in the current directory
    $ apex zip > /tmp/out.zip
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fn := &function.Function{
		Service: lambda.New(session.New(aws.NewConfig())),
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
	case args["invoke"].(bool):
		invoke(fn, args["--verbose"].(bool), args["--async"].(bool))
	case args["zip"].(bool):
		zip(fn)
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
		log.Fatalf("error: %s", err)
	}
}

// zip outputs the generated archive to stdout.
func zip(fn *function.Function) {
	zip, err := fn.Zip()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	_, err = io.Copy(os.Stdout, zip)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
