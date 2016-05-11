package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/stats"

	// commands
	_ "github.com/apex/apex/cmd/apex/autocomplete"
	_ "github.com/apex/apex/cmd/apex/build"
	_ "github.com/apex/apex/cmd/apex/delete"
	_ "github.com/apex/apex/cmd/apex/deploy"
	_ "github.com/apex/apex/cmd/apex/docs"
	_ "github.com/apex/apex/cmd/apex/infra"
	_ "github.com/apex/apex/cmd/apex/init"
	_ "github.com/apex/apex/cmd/apex/invoke"
	_ "github.com/apex/apex/cmd/apex/list"
	_ "github.com/apex/apex/cmd/apex/logs"
	_ "github.com/apex/apex/cmd/apex/metrics"
	_ "github.com/apex/apex/cmd/apex/rollback"
	_ "github.com/apex/apex/cmd/apex/upgrade"
	_ "github.com/apex/apex/cmd/apex/version"

	// plugins
	_ "github.com/apex/apex/plugins/env"
	_ "github.com/apex/apex/plugins/golang"
	_ "github.com/apex/apex/plugins/hooks"
	_ "github.com/apex/apex/plugins/inference"
	_ "github.com/apex/apex/plugins/java"
	_ "github.com/apex/apex/plugins/nodejs"
	_ "github.com/apex/apex/plugins/python"
	_ "github.com/apex/apex/plugins/shim"
)

// Terraform commands.
var tf = []string{
	"apply",
	"destroy",
	"get",
	"graph",
	"init",
	"output",
	"plan",
	"refresh",
	"remote",
	"show",
	"taint",
	"untaint",
	"validate",
	"version",
}

func main() {
	log.SetHandler(cli.Default)

	args := os.Args[1:]

	// Cobra does not (currently) allow us to pass flags for a sub-command
	// as if they were arguments, so we inject -- here after the first TF command.
	// TODO(tj): replace with a real solution and send PR to Cobra #251
	if len(os.Args) > 1 && os.Args[1] == "infra" {
		off := 1

	out:
		for i, a := range args {
			for _, cmd := range tf {
				if a == cmd {
					off = i
					break out
				}
			}
		}

		args = append(args[0:off], append([]string{"--"}, args[off:]...)...)
	}

	root.Command.SetArgs(args)

	if err := root.Command.Execute(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	stats.Client.ConditionalFlush(500, 24*time.Hour)
}
