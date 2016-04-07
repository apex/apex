package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apex/apex/boot/boilerplate"
	"github.com/apex/apex/infra"
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
	"defaultEnvironment": "dev",
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

var remoteStateCommand = `
  terraform remote config \
    -backend=s3 \
    -backend-config="region=%s" \
    -backend-config="bucket=%s" \
    -backend-config="key=terraform/state/%s"
`

// All bootstraps a project.
func All(region string) error {
	help("Enter the name of your project. It should be machine-friendly, as this\nis used to prefix your functions in Lambda.")
	name := prompt.StringRequired(indent("  Project name: "))

	help("Enter an optional description of your project.")
	description := prompt.String(indent("  Project description: "))

	fmt.Println()
	if prompt.Confirm(indent("Would you like to manage infrastructure with Terraform? (yes/no) ")) {
		fmt.Println()
		if err := initProject(name, description, ""); err != nil {
			return err
		}

		help("List the environments you would like (comma separated, e.g.: 'dev, test, prod')")
		envs := readEnvs(prompt.String(indent("  Environments: ")))

		fmt.Println()
		if err := initInfra(envs); err != nil {
			return err
		}

		fmt.Println()
		if prompt.Confirm(indent("Would you like to store Terraform state on S3? (yes/no) ")) {
			help("Enter the S3 bucket name for managing Terraform state (bucket needs\nto exist, use separate bucket for each project).")
			bucket := prompt.StringRequired(indent("  S3 bucket name: "))
			fmt.Println()

			if err := setupRemoteState(region, bucket, envs); err != nil {
				return err
			}
		}

		help("Setup complete!\n\nNext steps: \n  - apex infra plan - show an execution plan for Terraform configs\n  - apex infra apply - apply Terraform configs\n  - apex deploy - deploy example function")
		return nil
	}

	help("Enter IAM role used by Lambda functions.")
	iamRole := prompt.StringRequired(indent("  IAM role: "))

	fmt.Println()
	if err := initProject(name, description, iamRole); err != nil {
		return err
	}

	help("Setup complete!\n\nNext step: \n  - apex deploy - deploy example function")
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
func initInfra(envs []string) error {
	if _, err := exec.LookPath("terraform"); err != nil {
		return fmt.Errorf("terraform is not installed")
	}

	logf("creating ./infrastructure")

	if err := boilerplate.RestoreAssets(".", filepath.Join(infra.Dir, "modules")); err != nil {
		return err
	}

	for _, env := range envs {
		if err := setupEnv(env); err != nil {
			return err
		}

		if err := setupModules(env); err != nil {
			return err
		}
	}

	return nil
}

// setupEnv creates environment dir
func setupEnv(env string) error {
	logf("creating %s environment", env)

	if err := os.MkdirAll(filepath.Join(infra.Dir, env), 0755); err != nil {
		return err
	}

	maintf := filepath.Join(infra.Dir, "_env", "main.tf")
	data, err := boilerplate.Asset(maintf)
	if err != nil {
		return err
	}
	info, err := boilerplate.AssetInfo(maintf)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(infra.Dir, env, "main.tf"), data, info.Mode()); err != nil {
		return err
	}

	return nil
}

// setupModules performs a `terraform get`.
func setupModules(env string) error {
	logf("fetching %s modules", env)
	dir := filepath.Join(infra.Dir, env)
	return shell(modulesCommand, dir)
}

// setupRemoteState performs a `terraform remote config`.
func setupRemoteState(region, bucket string, envs []string) error {
	for _, env := range envs {
		logf("setting up remote %s state in bucket %q", env, bucket)
		cmd := fmt.Sprintf(remoteStateCommand, region, bucket, env)
		dir := infra.Dir
		if err := shell(cmd, dir); err != nil {
			return err
		}
	}

	return nil
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

// readEnvs splits string and removes whitespaces
func readEnvs(env string) (envs []string) {
	for _, e := range strings.Split(env, ",") {
		e := strings.TrimSpace(e)
		if e != "" {
			envs = append(envs, e)
		}
	}

	return
}

// help string output.
func help(s string) {
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(indent(s))
	os.Stdout.WriteString("\n\n")
}

// indent multiline string with 2 spaces.
func indent(s string) (out string) {
	for _, l := range strings.SplitAfter(s, "\n") {
		out += fmt.Sprintf("  %s", l)
	}

	return
}

// logf outputs a log message.
func logf(s string, v ...interface{}) {
	fmt.Printf("  \033[34m[+]\033[0m %s\n", fmt.Sprintf(s, v...))
}
