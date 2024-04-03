package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkCommand = &cobra.Command {
  Use: "link",
  Run: runLinkCommand,
  Short: "generate symlinks according to config",
  Long: "add longer description", // TODO add longer description here
}

func init() {
  RootCmd.AddCommand(linkCommand)
}

func runLinkCommand(cmd *cobra.Command, args []string) {
  fs := FileSystem
  fmt.Println("Symlinking dotfiles...")
  dotfileRoot := viper.Get("dotfile-path").(string)

  links := viper.GetStringMapString("links")

  for configName, configPath := range links {
    if configName == ".git"  || configName == "dotctl" {
      continue
    }
    dotPath := filepath.Join(dotfileRoot, configName)

    if configPath == ""{
      fmt.Fprintf(cmd.OutOrStdout(), "Warning: could not find config for %s\n", configName)
    }


    // destination needs to be removed before symlink
    if(DryRun) {
      log.Printf("Existing directory %s will be removed\n", configPath)

    } else {
      fs.RemoveAll(configPath)
    }

    testing := viper.Get("testing")

    if(DryRun) {
      log.Printf("Will link %s -> %s\n", configPath, dotPath)
    } else {
      if(testing == true) {
        fmt.Fprintf(cmd.OutOrStdout(), "%s,%s", configPath, dotPath)
      } else {
        err := afero.OsFs.SymlinkIfPossible(afero.OsFs{}, dotPath, configPath)
        if err != nil {
          log.Fatalf("Cannot symlink %s: %s\n", configName, err.Error())
        } else {
          fmt.Printf("%s linked\n", configName)
        }
      }
    }
  }


}
