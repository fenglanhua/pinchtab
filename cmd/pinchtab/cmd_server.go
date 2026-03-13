package main

import (
	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/pinchtab/pinchtab/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start full server",
	Run: func(cmd *cobra.Command, args []string) {
		maybeRunWizard()
		cfg := config.Load()
		server.RunDashboard(cfg, version)
	},
}

func init() {
	serverCmd.GroupID = "primary"
	rootCmd.AddCommand(serverCmd)
}
