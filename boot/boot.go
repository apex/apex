package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/apex/apex/boot/boilerplate"
	"github.com/apex/apex/infra"
	"github.com/mitchellh/go-wordwrap"
	"github.com/tj/go-prompt"
)

var modulesCommand = `
  terraform get
`

var projectConfig = `
{
  "name": "%s",
  "description": "%s",
  "memory": 128,
  "timeout": 5,
  "role": "%s",
  "environment": {}
}`

var projectConfigWithoutRole = `
{
  "name": "%s",
  "description": "%s",
  "memory": 128,
  "timeout": 5,
  "environment": {}
}`

// All bootstraps a project.
func All() error {
	help(`Enter the name of your project. It should be machine-friendly, as this is used to prefix your functions in Lambda.`)
	name := prompt.StringRequired("  Project name: ")

	help(`Enter an optional description of your project.`)
	description := prompt.String("  Project description: ")

	fmt.Println()
	if prompt.Confirm("Would you like to manage infrastructure with Terraform? (yes/no) ") {
		fmt.Println()
		if err := initProject(name, description, ""); err != nil {
			return err
		}

		if err := initInfra(); err != nil {
			return err
		}

		help("Setup complete!\n\nNext steps: \n  - apex infra plan - show an execution plan for Terraform configs\n  - apex infra apply - apply Terraform configs\n  - apex deploy - deploy example function")
	} else {
		fmt.Println()
		help(`Enter IAM role used by Lambda functions.`)
		iamRole := prompt.StringRequired("  IAM role: ")

		fmt.Println()
		if err := initProject(name, description, iamRole); err != nil {
			return err
		}

		help("Setup complete!\n\nNext step: \n  - apex deploy - deploy example function")
	}

	return nil
}

// Project bootstraps a project.
func initProject(name, description, iamRole string) error {
	logf("creating ./project.json")

	var project string
	if iamRole == "" {
		project = fmt.Sprintf(projectConfigWithoutRole, name, description)
	} else {
		project = fmt.Sprintf(projectConfig, name, description, iamRole)
	}

	if err := ioutil.WriteFile("project.json", []byte(project), 0644); err != nil {
		return err
	}

	logf("creating ./functions")
	return boilerplate.RestoreAssets(".", "functions")
}

// infra bootstraps terraform for infrastructure management.
func initInfra() error {
	if _, err := exec.LookPath("terraform"); err != nil {
		return fmt.Errorf("terraform is not installed")
	}

	logf("creating ./infrastructure")
	if err := boilerplate.RestoreAssets(".", infra.Dir); err != nil {
		return err
	}

	return setupModules()
}

// setupModules performs a `terraform get`.
func setupModules() error {
	logf("fetching modules")
	return shell(modulesCommand, infra.Dir)
}

// shell executes `command` in a shell within `dir`.
func shell(command, dir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command: %s: %s", out, err)
	}

	return nil
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
