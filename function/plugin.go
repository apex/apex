package function

import (
	"fmt"
)

// defaultPlugins are the default plugins which are required by Apex. Note that
// the order here is important for some plugins such as inference before the
// runtimes.
var defaultPlugins = []string{
	"inference",
	"golang",
	"python",
	"nodejs",
	"hooks",
	"env",
	"shim",
}

// Hook type.
type Hook string

// Hooks available.
const (
	// OpenHook is called when the function configuration is first loaded.
	OpenHook Hook = "open"

	// BuildHook is called when a build is started.
	BuildHook = "build"

	// CleanHook is called when a build is complete.
	CleanHook = "clean"

	// DeployHook is called after build and before a deploy.
	DeployHook = "deploy"
)

// A Plugin is a chunk of isolated(ish) logic which reacts to various
// hooks within the system in order to implement specific features
// such as runtime inference or environment variable support.
type Plugin interface {
	Run(hook Hook, fn *Function) error
}

// Registered plugins.
var plugins = make(map[string]Plugin)

// Register plugin by `name`.
func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}

// ByName returns the plugin by `name`.
func ByName(name string) (Plugin, error) {
	if v, ok := plugins[name]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("invalid plugin %q", name)
	}
}

// hook runs the default plugins, and those defined by Plugins in sequence.
func (f *Function) hook(hook Hook) error {
	for _, name := range defaultPlugins {
		plugin, err := ByName(name)
		if err != nil {
			return err
		}

		if err := plugin.Run(hook, f); err != nil {
			return fmt.Errorf("%s: %s", name, err)
		}
	}

	return nil
}
