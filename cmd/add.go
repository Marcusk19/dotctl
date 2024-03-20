package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Marcusk19/bender/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
  RootCmd.AddCommand(addCommand)
}

var addCommand = &cobra.Command {
  Use: "add",
  Short: "Adds config to be tracked by bender",
  Long: "TODO: add longer description", // TODO add more description
  Run: runAddCommand,
}

func runAddCommand(cmd *cobra.Command, args []string) {
  fs := FileSystem

  if len(args) <= 0 {
    fmt.Println("ERROR: requires at least one argument")
    return
  }

  configSrc := args[0]
  dirs := strings.Split(configSrc, "/")
  name := dirs[len(dirs) - 1]
  viper.Set(name, configSrc)
  viper.WriteConfig()

  dotfileDest := filepath.Join(DotfilePath, name)

  if DryRun {
    fmt.Printf("Will copy %s -> %s \n", configSrc, dotfileDest)
    return
  }
  
  err := tools.CopyDir(fs, configSrc, dotfileDest)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Copied %s -> %s\n", configSrc, dotfileDest)
}
