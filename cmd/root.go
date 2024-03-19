/*
Copyright © 2024 Marcus Kok
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
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

var FileSystem afero.Fs

func init() {
  // define flags and config sections


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
  viper.BindPFlag("config-path", RootCmd.PersistentFlags().Lookup("config-path"))

  viper.BindEnv("testing")
  viper.SetDefault("testing", false)

  viper.SetConfigName("bender.yml")
  viper.SetConfigType("yaml")
  viper.AddConfigPath(filepath.Join(defaultDotPath, "bender"))
  viper.AddConfigPath("./bender")

  err := viper.ReadInConfig()
  if err != nil {
    fmt.Println("No config detected. You can generate one by using 'bender init'")
  }

  FileSystem = UseFilesystem()

}

func UseFilesystem() afero.Fs {
  testing := viper.Get("testing")
  if(testing == "true") {
    return afero.NewMemMapFs()
  } else {
    return afero.NewOsFs()
  }
}

// TODO: this can probably be removed
func SetUpForTesting() afero.Fs {
  viper.Set("testing", true)
  fs := UseFilesystem()

  homedir := "bender_test/"
  fs.MkdirAll(filepath.Join(homedir, ".config/"), 0755)
  fs.MkdirAll(filepath.Join(homedir, ".dotfiles/"), 0755)

  fs.Mkdir("bin/", 0755)
  fs.Create("bin/alacritty")
  fs.Create("bin/nvim")
  fs.Create("bin/tmux")

  fs.Mkdir(filepath.Join(homedir, ".config/alacritty"), 0755)
  fs.Mkdir(filepath.Join(homedir, ".config/nvim"), 0755)
  fs.Mkdir(filepath.Join(homedir, ".config/tmux"), 0755)

  fs.Create(filepath.Join(homedir, ".config/alacritty/alacritty.conf"))
  fs.Create(filepath.Join(homedir, ".config/nvim/nvim.conf"))
  fs.Create(filepath.Join(homedir, ".config/tmux/tmux.conf"))

  FileSystem = fs
 
  return fs
}

