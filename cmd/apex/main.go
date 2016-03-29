package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/apex/apex/cmd/apex/root"

	// commands
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

func main() {
	log.SetHandler(cli.Default)

	if err := root.Command.Execute(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
