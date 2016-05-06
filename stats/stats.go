package stats

import (
	"runtime"

	"github.com/tj/go-cli-analytics"

	"github.com/apex/apex/cmd/apex/version"
)

// Client for Segment analytics.
var Client = analytics.New(&analytics.Config{
	WriteKey: "AgKH1g5KH9cJWcXw5djwsurxGcQWfHS6",
	Dir:      ".apex",
})

// Track event `name` with optional `props`.
func Track(name string, props map[string]interface{}) {
	if props == nil {
		props = map[string]interface{}{}
	}

	props["version"] = version.Version
	props["os"] = runtime.GOOS
	props["arch"] = runtime.GOARCH

	Client.Track(name, props)
}
