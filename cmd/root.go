package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "git-release-manager",
		Short: "Git release manager",
		Long:  "GRM is a Git release manager.",
	}
)

// Execute the root command
func Execute() error {
	return rootCmd.Execute()
}
