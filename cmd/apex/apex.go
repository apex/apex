package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/apex/apex/runtime/golang"
	_ "github.com/apex/apex/runtime/nodejs"
	_ "github.com/apex/apex/runtime/python"

	"github.com/apex/apex/function"
	"github.com/apex/apex/logs"
	"github.com/apex/apex/project"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/segmentio/go-prompt"
	"github.com/tj/docopt"
)

var version = "0.3.0"

const usage = `
  Usage:
    apex deploy [options] [<name>...]
    apex delete [options] [<name>...]
    apex invoke [options] <name> [--async] [-v]
    apex rollback [options] <name> [<version>]
    apex logs [options] <name> [--filter pattern]
    apex build [options] <name>
    apex list [options]
    apex -h | --help
    apex --version

  Options:
    -F, --filter pattern    Filter logs with pattern [default: ]
    -l, --log-level level   Log severity level [default: info]
    -a, --async             Async invocation
    -C, --chdir path        Working directory
    -y, --yes               Automatic yes to prompts
    -h, --help              Output help information
    -v, --verbose           Output verbose logs
    -V, --version           Output version

  Examples:
    Deploy all functions
    $ apex deploy

    Deploy specific functions
    $ apex deploy foo bar

    Delete all functions
    $ apex delete

    Delete specified functions
    $ apex delete foo bar

    Invoke a function with input json
    $ apex invoke foo < request.json

    Rollback a function to the previous version
    $ apex rollback foo

    Rollback a function to the specified version
    $ apex rollback bar 3

    Deploy functions in a different project
    $ apex deploy -C ~/dev/myapp

    Build zip output for a function
    $ apex build foo > /tmp/out.zip
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	log.SetHandler(text.New(os.Stderr))

	if l, err := log.ParseLevel(args["--log-level"].(string)); err == nil {
		log.SetLevel(l)
	}

	project := &project.Project{
		Service: lambda.New(session.New(aws.NewConfig())),
		Log:     log.Log,
		Path:    ".",
	}

	if dir, ok := args["--chdir"].(string); ok {
		if err := os.Chdir(dir); err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	if err := project.Open(); err != nil {
		log.Fatalf("error opening project: %s", err)
	}

	switch {
	case args["list"].(bool):
		list(project)
	case args["deploy"].(bool):
		deploy(project, args["<name>"].([]string))
	case args["delete"].(bool):
		delete(project, args["<name>"].([]string), args["--yes"].(bool))
	case args["invoke"].(bool):
		invoke(project, args["<name>"].([]string), args["--verbose"].(bool), args["--async"].(bool))
	case args["rollback"].(bool):
		rollback(project, args["<name>"].([]string), args["<version>"])
	case args["build"].(bool):
		build(project, args["<name>"].([]string))
	case args["logs"].(bool):
		tail(project, args["<name>"].([]string), args["--filter"].(string))
	}
}

// list functions.
func list(project *project.Project) {
	// TODO(tj): more informative output
	fmt.Println()
	for _, fn := range project.Functions {
		fmt.Printf("  - %s (%s)\n", fn.Name, fn.Runtime)
	}
	fmt.Println()
}

// invoke reads request json from stdin and outputs the responses.
func invoke(project *project.Project, name []string, verbose, async bool) {
	dec := json.NewDecoder(os.Stdin)
	kind := function.RequestResponse

	if async {
		kind = function.Event
	}

	fn, err := project.FunctionByName(name[0])
	if err != nil {
		log.Fatalf("error: %s", err)
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

		// TODO(tj) rename flag to --with-logs or --logs
		if verbose {
			io.Copy(os.Stderr, logs)
		}

		io.Copy(os.Stdout, reply)
		fmt.Fprintf(os.Stdout, "\n")
	}
}

// deploy code and config changes.
func deploy(project *project.Project, names []string) {
	var err error

	if len(names) == 0 {
		names, err = project.FunctionNames()
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if err := project.DeployAndClean(names); err != nil {
		log.Fatalf("error: %s", err)
	}
}

// delete the functions.
func delete(project *project.Project, names []string, force bool) {
	var err error

	if len(names) == 0 {
		names, err = project.FunctionNames()
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if !force && len(names) > 1 {
		fmt.Printf("The following will be deleted:\n\n")
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Printf("\n")
	}

	if !force && !prompt.Confirm("Are you sure? (yes/no)") {
		return
	}

	if err := project.Delete(names); err != nil {
		log.Fatalf("error: %s", err)
	}
}

// rollback the function with optional version.
func rollback(project *project.Project, name []string, version interface{}) {
	fn, err := project.FunctionByName(name[0])
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if version == nil {
		err = fn.Rollback()
	} else {
		err = fn.RollbackVersion(version.(string))
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// build outputs the generated archive to stdout.
func build(project *project.Project, name []string) {
	fn, err := project.FunctionByName(name[0])
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	zip, err := fn.Zip()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	_, err = io.Copy(os.Stdout, zip)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// tail outputs logs with optional filter pattern.
func tail(project *project.Project, name []string, filter string) {
	service := cloudwatchlogs.New(session.New(aws.NewConfig()))

	// TODO(tj): refactor logs.Logs to take Project so this hack
	// can be removed, it'll also make multi-function tailing easier
	group := fmt.Sprintf("/aws/lambda/%s_%s", project.Name, name[0])

	l := logs.Logs{
		LogGroupName:  group,
		FilterPattern: filter,
		Service:       service,
		Log:           log.Log,
	}

	for event := range l.Tail() {
		fmt.Printf("%s", *event.Message)
	}

	if err := l.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
