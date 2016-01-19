package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	_ "github.com/apex/apex/runtime/golang"
	_ "github.com/apex/apex/runtime/nodejs"
	_ "github.com/apex/apex/runtime/python"
)

const version = "0.4.1"

func main() {
	log.SetHandler(cli.Default)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
