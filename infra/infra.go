// Package infra proxies Terraform commands.
package infra

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/apex/apex/function"
	"github.com/apex/log"
)

// TODO(tj): revisit removal of Region when TF supports ~/.aws/config

// Dir in which Terraform configs are stored
const Dir = "infrastructure"

// Proxy is a wrapper around Terraform commands.
type Proxy struct {
	Functions []*function.Function
	Region    string
}

// Run terraform command in infrastructure directory.
func (p *Proxy) Run(args ...string) error {
	if p.shouldInjectVars(args) {
		args = append(args, p.functionVars()...)
	}

	log.WithFields(log.Fields{
		"args": args,
	}).Debug("terraform")

	cmd := exec.Command("terraform", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = Dir

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("AWS_REGION=%s", p.Region))

	return cmd.Run()
}

// functionVars returns the function ARN's as terraform -var arguments.
func (p *Proxy) functionVars() (args []string) {
	for _, fn := range p.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		args = append(args, "-var")
		args = append(args, fmt.Sprintf("apex_function_%s=%s", fn.Name, *config.Configuration.FunctionArn))
	}

	return args
}

// shouldInjectVars checks if the command accepts -var flags.
func (p *Proxy) shouldInjectVars(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return args[0] == "plan" || args[0] == "apply"
}

// Output fetches output variable `name` from terraform.
func Output(name string) (string, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("terraform output %s", name))
	cmd.Dir = Dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}
