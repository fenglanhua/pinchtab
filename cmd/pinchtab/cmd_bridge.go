package main

import (
	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/pinchtab/pinchtab/internal/server"
	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Start single-instance bridge-only server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		server.RunBridgeServer(cfg)
	},
}

func init() {
	bridgeCmd.GroupID = "primary"
	rootCmd.AddCommand(bridgeCmd)
}
