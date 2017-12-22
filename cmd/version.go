// Copyright © 2017 Mladen Popadic <mladen.popadic.4@gmail.com>

package cmd

import (
	"github.com/mpopadic/go_n_find/colors"
	"github.com/spf13/cobra"
)

// Version of CLI tool
var (
	Version = "No version provided"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version number",
	Long:  `Prints version number`,
	Run: func(cmd *cobra.Command, args []string) {
		colors.CYAN.Printf("version: %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
