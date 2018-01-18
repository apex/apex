// Package upgrade fetches the latest Apex binary, if any, for the current platform.
package upgrade

import (
	"runtime"

	"github.com/pkg/errors"
	"github.com/tj/go-update"
	"github.com/tj/go-update/stores/github"
	"github.com/tj/go/term"

	"github.com/apex/apex/internal/progressreader"
	"github.com/apex/apex/internal/util"
)

// Upgrade the current `version` of apex to the latest.
func Upgrade(version string) error {
	term.HideCursor()
	defer term.ShowCursor()

	p := &update.Manager{
		Command: "apex",
		Store: &github.Store{
			Owner:   "apex",
			Repo:    "apex",
			Version: version,
		},
	}

	// fetch releases
	releases, err := p.LatestReleases()
	if err != nil {
		return errors.Wrap(err, "fetching latest release")
	}

	// no updates
	if len(releases) == 0 {
		util.LogPad("No updates available, you're good :)")
		return nil
	}

	// latest
	r := releases[0]

	// find the tarball for this system
	a := r.FindTarball(runtime.GOOS, runtime.GOARCH)
	if a == nil {
		return errors.Errorf("failed to find a binary for %s %s", runtime.GOOS, runtime.GOARCH)
	}

	// download tarball to a tmp dir
	tarball, err := a.DownloadProxy(progressreader.New)
	if err != nil {
		return errors.Wrap(err, "downloading tarball")
	}

	// install it
	if err := p.Install(tarball); err != nil {
		return errors.Wrap(err, "installing")
	}

	term.ClearAll()
	util.LogPad("Updated %s to %s", version, r.Version)

	return nil
}
