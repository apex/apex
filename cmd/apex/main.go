package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	_ "github.com/apex/apex/plugins/env"
	_ "github.com/apex/apex/plugins/golang"
	_ "github.com/apex/apex/plugins/hooks"
	_ "github.com/apex/apex/plugins/inference"
	_ "github.com/apex/apex/plugins/nodejs"
	_ "github.com/apex/apex/plugins/python"
	_ "github.com/apex/apex/plugins/shim"
)

const version = "0.5.0"

func main() {
	log.SetHandler(cli.Default)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
