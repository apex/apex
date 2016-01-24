// Package hooks implements hook script support.
package hooks

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/apex/function"
	"github.com/apex/log"
)

func init() {
	function.RegisterPlugin("hooks", &Plugin{})
}

// A HookError represents a failed hook command.
type HookError struct {
	Hook    function.Hook
	Command string
	Output  string
}

// Error string.
func (e *HookError) Error() string {
	return fmt.Sprintf("%s: %s", e.Hook, e.Output)
}

// Plugin implementation.
type Plugin struct{}

// Run executes any commands defined for a hook.
func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	var command string

	switch hook {
	case function.CleanHook:
		command = fn.Hooks.Clean
	case function.BuildHook:
		command = fn.Hooks.Build
	case function.DeployHook:
		command = fn.Hooks.Deploy
	}

	if command == "" {
		return nil
	}

	fn.Log.WithFields(log.Fields{
		"hook":    hook,
		"command": command,
	}).Debug("hook")

	cmd := exec.Command("sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("FUNCTION=%s", fn.Name))
	cmd.Dir = fn.Path

	b, err := cmd.CombinedOutput()
	if err != nil {
		return &HookError{
			Hook:    hook,
			Command: command,
			Output:  string(b),
		}
	}

	return nil
}
