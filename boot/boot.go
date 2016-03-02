package boot

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/apex/boot/boilerplate"
	"github.com/mitchellh/go-wordwrap"
	"github.com/tj/go-prompt"
)

var projectConfig = `
{
  "name": "%s",
  "description": "%s",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::%s:role/lambda",
  "environment": {}
}`

// All bootstraps a project.
func All() error {
	if err := Project(); err != nil {
		return err
	}

	help(`Setup complete :)`)

	return nil
}

// Project bootstraps a project.
func Project() error {
	help(`Enter the name of your project. It should be machine-friendly, as this is used to prefix your functions in Lambda.`)
	name := prompt.StringRequired("  Project name: ")

	help(`Enter an optional description of your project.`)
	description := prompt.String("  Project description: ")

	help(`Enter your AWS Account ID.`)
	accountID := prompt.StringRequired("  AWS Account ID: ")
	fmt.Println()

	logf("creating ./project.json")
	project := fmt.Sprintf(projectConfig, name, description, accountID)
	if err := ioutil.WriteFile("project.json", []byte(project), 0644); err != nil {
		return err
	}

	logf("creating ./functions")
	return boilerplate.RestoreAssets(".", "functions")
}

// help string output.
func help(s string) {
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(wordwrap.WrapString(s, 70))
	os.Stdout.WriteString("\n\n")
}

// logf outputs a log message.
func logf(s string, v ...interface{}) {
	fmt.Printf("  \033[34m[+]\033[0m %s\n", fmt.Sprintf(s, v...))
}
