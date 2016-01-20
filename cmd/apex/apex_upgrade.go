package main

import (
	"github.com/apex/apex/upgrade"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:              "upgrade",
	Short:            "Ugrade apex to the latest stable release",
	PersistentPreRun: pv.noopRun,
	Run:              upgradeCmdRun,
}

func upgradeCmdRun(c *cobra.Command, args []string) {
	err := upgrade.Upgrade(version)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
