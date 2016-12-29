// Package exec proxies all commands.
package exec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/apex/function"
	"github.com/apex/log"
)

// Proxy is a wrapper around commands.
type Proxy struct {
	Functions   []*function.Function
	Environment string
	Region      string
	Role        string
	Dir         string
}

// Run command in specified directory.
func (p *Proxy) Run(command string, args ...string) error {
	log.WithFields(log.Fields{
		"command": command,
		"args":    args,
	}).Debug("exec")

	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), p.functionEnvVars()...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = p.Dir

	return cmd.Run()
}

func (p *Proxy) functionEnvVars() (args []string) {
	args = append(args, fmt.Sprintf("apex_environment=%s", p.Environment))

	if p.Role != "" {
		args = append(args, fmt.Sprintf("apex_function_role=%s", p.Role))
	}

	for _, fn := range p.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		args = append(args, fmt.Sprintf("apex_function_%s=%s", fn.Name, *config.Configuration.FunctionArn))
	}

	return args
}
