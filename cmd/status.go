package cmd

import (
	"fmt"
	"log"
	"slices"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
  RootCmd.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command {
  Use: "status",
  Short: "View status of dotctl",
  Long: "TODO: add longer description",
  Run: runStatusCommand,
}

func runStatusCommand(cmd *cobra.Command, args[]string) {
  fs := FileSystem
  links := viper.GetStringMapString("links")

  var ignoredDirs = []string{".git", "dotctl", ".gitignore"}
  
  dotfiles, err := afero.ReadDir(fs, viper.GetString("dotfile-path"))
  if err != nil {
    log.Fatalf("Cannot read dotfile dir: %s\n", err)
  }

  var linkedConfigs []string
  var orphanedConfigs []string

  fmt.Fprintln(cmd.OutOrStdout(), "Config directories currently in dotfile path:\n")
  for _, dotfileDir := range(dotfiles) {
    dirName := dotfileDir.Name()
    if !slices.Contains(ignoredDirs, dirName) {
      if links[dirName] != "" {
        // fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s\n", dirName, links[dirName]) 
        linkedConfigs = append(linkedConfigs, dirName, links[dirName])
      } else {
        // fmt.Fprintln(cmd.OutOrStdout(), dirName)
        orphanedConfigs = append(orphanedConfigs, dirName)
      }
    }
  }

  for i := 0; i < len(linkedConfigs); i += 2 {
    fmt.Fprintf(cmd.OutOrStdout(), "%s (links to %s)\n", linkedConfigs[i], linkedConfigs[i+1])
  }
  fmt.Fprintln(cmd.OutOrStdout(), "================")
  
  fmt.Fprintln(cmd.OutOrStdout(), "Orphaned configs")

  for _, conf := range(orphanedConfigs) {
    fmt.Fprintln(cmd.OutOrStdout(), conf)
  }


    
}


