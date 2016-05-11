package function

import "github.com/apex/apex/archive"

// A Plugin is a chunk of isolated(ish) logic which reacts to various
// hooks within the system in order to implement specific features
// such as runtime inference or environment variable support.
type Plugin interface{}

// Opener reacts to the Open hook.
type Opener interface {
	Open(*Function) error
}

// Builder reacts to the Build hook.
type Builder interface {
	Build(*Function, *archive.Zip) error
}

// Cleaner reacts to the Clean hook.
type Cleaner interface {
	Clean(*Function) error
}

// Deployer reacts to the Deploy hook.
type Deployer interface {
	Deploy(*Function) error
}

// Registered plugins.
var plugins = make(map[string]Plugin)

// RegisterPlugin registers `plugin` by `name`.
func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}
