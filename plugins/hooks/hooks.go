// Package hooks implements hook script support.
package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/apex/log"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("hooks", &Plugin{})
}

// HookError represents a failed hook command.
type HookError struct {
	Hook    string
	Command string
	Output  string
}

// Error string.
func (e *HookError) Error() string {
	return fmt.Sprintf("%s hook: %s", e.Hook, e.Output)
}

// Plugin implementation.
type Plugin struct{}

// Build runs the "build" hook commands.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	return p.run("build", fn.Hooks.Build, fn)
}

// Clean runs the "clean" hook commands.
func (p *Plugin) Clean(fn *function.Function) error {
	return p.run("clean", fn.Hooks.Clean, fn)
}

// Deploy runs the "deploy" hook commands.
func (p *Plugin) Deploy(fn *function.Function) error {
	return p.run("deploy", fn.Hooks.Deploy, fn)
}

// run a hook command.
func (p *Plugin) run(hook, command string, fn *function.Function) error {
	if command == "" {
		return nil
	}

	fn.Log.WithFields(log.Fields{
		"hook":    hook,
		"command": command,
	}).Debug("hook")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Env = os.Environ()
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
