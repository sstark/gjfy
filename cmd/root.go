package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version string
)

var rootCmd = &cobra.Command{
	Use:   "gjfy",
	Short: "gjfy one-time links",
	Version: Version,
	Long:  `
gjfy is a web service and tool for creating and providing one-time clickable links`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
