package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/apex/apex/function"
	"github.com/apex/apex/project"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/segmentio/go-prompt"
	"github.com/tj/docopt"
)

var version = "0.1.0"

const usage = `
  Usage:
    apex deploy [options] [<name>...]
    apex delete [options] [<name>...]
    apex invoke [options] <name> [--async] [-v]
    apex rollback [options] <name> [<version>]
    apex build [options] <name>
    apex list [options]
    apex -h | --help
    apex --version

  Options:
    -l, --log-level level   Log severity level [default: info]
    -a, --async             Async invocation
    -C, --chdir path        Working directory
    -y, --yes               Automatic yes to prompts
    -h, --help              Output help information
    -v, --verbose           Output verbose logs
    -V, --version           Output version

  Examples:
    Deploy a function in the current directory
    $ apex deploy

    Delete a function in the current directory
    $ apex delete

    Invoke a function in the current directory
    $ apex invoke < request.json

    Rollback a function to the previous version (another call will go back to the latest version)
    $ apex rollback

    Rollback a function to the specified version (another call will go back to the latest version)
    $ apex rollback 3

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

// deploy changes.
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

// delete the resources.
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

// rollback the function.
func rollback(project *project.Project, name []string, version interface{}) {
	fn, err := project.FunctionByName(name[0])
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if version == nil {
		err = fn.Rollback()
	} else {
		err = fn.Rollback(version.(string))
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
