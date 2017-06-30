package exec

import (
	"errors"

	"github.com/apex/apex/exec"
	"github.com/tj/cobra"

	"os"

	"github.com/apex/apex/cmd/apex/root"
)

var dir string

// example output.
const example = `
    Run terraform apply command
    $ apex exec -d ./infrastructure/dev/pre terraform apply
`

// Command config.
var Command = &cobra.Command{
	Use:     "exec",
	Short:   "Command execution passthrough",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)
	f := Command.Flags()
	f.StringVarP(&dir, "dir", "d", "", "Which directory to execute command from")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	var command string

	if dir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		dir = pwd
	}

	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	} else {
		return errors.New("No command specified")
	}

	err := root.Project.LoadFunctions()
	if err != nil {
		return err
	}

	p := &exec.Proxy{
		Functions:   root.Project.Functions,
		Region:      *root.Session.Config.Region,
		Environment: root.Project.InfraEnvironment,
		Role:        root.Project.Role,
		Dir:         dir,
	}

	return p.Run(command, args...)
}
