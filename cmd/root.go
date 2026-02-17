package main

import (
	"github.com/audoctl/audoctl/cmd/audoctl"
	"github.com/audoctl/audoctl/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "audoctl",
	Short: "Audit and execution control for AI systems",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(
		version.VersionCmd,
		audoctl.CtlCmd,
	)
}
