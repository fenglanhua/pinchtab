package main

import (
	"fmt"
	"os"

	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/pinchtab/pinchtab/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP stdio server",
	Long:  "Start the Model Context Protocol stdio server and proxy browser actions to a running PinchTab instance.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runMCP(cfg)
	},
}

func init() {
	mcpCmd.GroupID = "primary"
	rootCmd.AddCommand(mcpCmd)
}

func runMCP(cfg *config.RuntimeConfig) {
	baseURL := os.Getenv("PINCHTAB_URL")
	if baseURL == "" {
		port := cfg.Port
		if port == "" {
			port = "9867"
		}
		baseURL = "http://127.0.0.1:" + port
	}

	token := os.Getenv("PINCHTAB_TOKEN")
	if token == "" {
		token = cfg.Token
	}

	mcp.Version = version

	if err := mcp.Serve(baseURL, token); err != nil {
		fmt.Fprintf(os.Stderr, "mcp server error: %v\n", err)
		os.Exit(1)
	}
}
