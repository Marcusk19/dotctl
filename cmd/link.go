package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var linkCommand = &cobra.Command {
  Use: "link",
  Run: runLinkCommand,
}

func init() {
  RootCmd.AddCommand(linkCommand)
}

func runLinkCommand(cmd *cobra.Command, args []string) {
  fs := UseFilesystem()
  fmt.Println("Symlinking dotfiles...")
  entries, err := afero.ReadDir(fs, DotfilePath)
  if err != nil {
    log.Fatal(err)
  }
  for _, entry := range(entries) {
    if entry.Name() == ".git" {
      continue
    }
    dotPath := filepath.Join(DotfilePath, entry.Name())
    configPath := filepath.Join(ConfigPath, entry.Name())

    // destination needs to be removed before symlink
    if(DryRun) {
      log.Printf("Existing directory %s will be removed\n", configPath)

    } else {
      fs.RemoveAll(configPath)
    }

    if(DryRun) {
      log.Printf("Will link %s -> %s\n", dotPath, configPath)
    } else {
      err = afero.OsFs.SymlinkIfPossible(afero.OsFs{}, dotPath, configPath)
    }
    if err != nil {
      log.Fatalf("Cannot symlink %s: %s", entry.Name(), err.Error())
    }
  }

}
