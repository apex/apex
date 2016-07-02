package boot

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/tj/go-prompt"

	"github.com/apex/apex/boot/boilerplate"
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

var projectConfig = `{
  "name": "%s",
  "description": "%s",
  "memory": 128,
  "timeout": 5,
  "role": "%s",
  "environment": {}
}`

var setupCompleteVanilla = `
  Setup complete, deploy those functions!

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
	return b.bootVanilla()
}

// check if there's a project.
func (b *Bootstrapper) isProject() bool {
	_, err := os.Stat("project.json")
	return err == nil
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
