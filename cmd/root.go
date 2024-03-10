/*
Copyright © 2024 Marcus Kok

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var RootCmd = &cobra.Command{
	Use:   "bender",
	Short: "dotfile management",
	Long: `Bender is a CLI tool for syncing your
  dotfiles. It provides an opiniated way to symlink
  a dotfile directory to various config directories.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var DotfilePath string
var ConfigPath string

func init() {
  // define flags and config sections

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bender.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
  defaultDotPath := os.Getenv("HOME") + "/.dotfiles/"
  defaultConfPath := os.Getenv("HOME") + "/.config/"
  RootCmd.PersistentFlags().StringVar(
    &DotfilePath,
    "dotfile-path",
    defaultDotPath,
    "Path pointing to dotfiles directory",
  )
  RootCmd.PersistentFlags().StringVar(
    &ConfigPath,
    "config-path",
    defaultConfPath,
    "Path pointing to config directory",
  )
  viper.BindPFlag("dotfile-path", RootCmd.PersistentFlags().Lookup("dotfile-path"))
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


