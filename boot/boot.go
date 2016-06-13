package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/tj/go-prompt"

	"github.com/apex/apex/boot/boilerplate"
	"github.com/apex/apex/infra"
)

// TODO(tj): attempt creation of S3 bucket to streamline that as well
// TODO(tj): idempotency, if the project exists then skip all this, or provision AWS
// TODO(tj): validate the name, \w+ or similar should be fine

var logo = `

             _    ____  _______  __
            / \  |  _ \| ____\ \/ /
           / _ \ | |_) |  _|  \  /
          / ___ \|  __/| |___ /  \
         /_/   \_\_|   |_____/_/\_\

`

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

var remoteStateCommand = `
  terraform remote config \
    -backend=s3 \
    -backend-config="region=%s" \
    -backend-config="bucket=%s" \
    -backend-config="key=terraform/state/%s"
`

var setupCompleteVanilla = `
  Setup complete, deploy those functions!

    $ apex deploy
`

var setupCompleteTerraform = `
  Setup complete, preview the infrastructure plan,
  apply it, then deploy those functions. Later you can
  change the environment with the --env flag.

    $ apex infra plan
    $ apex infra apply
    $ apex deploy
`

var iamAssumeRolePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

var iamLogsPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}`

// Bootstrapper initializes a project and AWS account for the user.
type Bootstrapper struct {
	IAM    iamiface.IAMAPI
	Region string

	name        string
	description string
}

// Boot the project.
func (b *Bootstrapper) Boot() error {
	fmt.Println(logo)

	if b.isProject() {
		help("I've detected a ./project.json file, this seems to already be a project!")
		return nil
	}

	help("Enter the name of your project. It should be machine-friendly, as this\nis used to prefix your functions in Lambda.")
	b.name = prompt.StringRequired(indent("  Project name: "))

	help("Enter an optional description of your project.")
	b.description = prompt.String(indent("  Project description: "))

	fmt.Println()
	if prompt.Confirm(indent("Would you like to manage infrastructure with Terraform? (yes/no) ")) {
		return b.bootTerraform()
	}

	fmt.Println()
	return b.bootVanilla()
}

// check if there's a project.
func (b *Bootstrapper) isProject() bool {
	_, err := os.Stat("project.json")
	return err == nil
}

// Bootstrap Terraform.
func (b *Bootstrapper) bootTerraform() error {
	fmt.Println()

	iamRole, err := b.createRole()
	if err != nil {
		return err
	}

	if err := b.initProjectFiles(iamRole); err != nil {
		return err
	}

	help("List the environments you would like (comma separated, e.g.: 'stage, prod')")
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

		if err := setupRemoteState(b.Region, bucket, envs); err != nil {
			return err
		}
	}

	fmt.Println(setupCompleteTerraform)
	return nil

}

// Bootstrap without Terraform.
func (b *Bootstrapper) bootVanilla() error {
	iamRole, err := b.createRole()
	if err != nil {
		return err
	}

	if err := b.initProjectFiles(iamRole); err != nil {
		return err
	}

	fmt.Println(setupCompleteVanilla)
	return nil
}

// Create IAM role, returning the ARN.
func (b *Bootstrapper) createRole() (string, error) {
	roleName := fmt.Sprintf("%s_lambda_function", b.name)
	policyName := fmt.Sprintf("%s_lambda_logs", b.name)

	logf("creating IAM %s role", roleName)
	role, err := b.IAM.CreateRole(&iam.CreateRoleInput{
		RoleName:                 &roleName,
		AssumeRolePolicyDocument: aws.String(iamAssumeRolePolicy),
	})

	if err != nil {
		return "", fmt.Errorf("creating role: %s", err)
	}

	logf("creating IAM %s policy", policyName)
	policy, err := b.IAM.CreatePolicy(&iam.CreatePolicyInput{
		PolicyName:     &policyName,
		Description:    aws.String("Allow lambda_function to utilize CloudWatchLogs. Created by apex(1)."),
		PolicyDocument: aws.String(iamLogsPolicy),
	})

	if err != nil {
		return "", fmt.Errorf("creating policy: %s", err)
	}

	logf("attaching policy to lambda_function role.")
	_, err = b.IAM.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  &roleName,
		PolicyArn: policy.Policy.Arn,
	})

	if err != nil {
		return "", fmt.Errorf("creating policy: %s", err)
	}

	return *role.Role.Arn, nil
}

// Initialize project files such as project.json and ./functions.
func (b *Bootstrapper) initProjectFiles(iamRole string) error {
	logf("creating ./project.json")

	project := fmt.Sprintf(projectConfig, b.name, b.description, iamRole)

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
		dir := filepath.Join(infra.Dir, env)
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
