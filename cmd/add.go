package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Marcusk19/dotctl/tools"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
  RootCmd.AddCommand(addCommand)
}

var addCommand = &cobra.Command {
  Use: "add",
  Short: "Adds config to be tracked by dotctl",
  Long: "TODO: add longer description", // TODO add more description
  Run: runAddCommand,
}

func runAddCommand(cmd *cobra.Command, args []string) {
  fs := FileSystem

  testing := viper.GetBool("testing")

  if len(args) <= 0 {
    fmt.Println("ERROR: requires config path")
    return
  }

  configSrc := args[0]
  dirs := strings.Split(configSrc, "/")
  name := dirs[len(dirs) - 1] // take the last section of the path, this should be the name

  links := viper.GetStringMap("links")
  links[name] = configSrc
  viper.Set("links", links)
  if !testing {
    err := viper.WriteConfig()
    if err != nil {
      fmt.Printf("Problem updating dotctl config %s", err)
    }
  }

  dotfilePath := viper.Get("dotfile-path").(string)

  dotfileDest := filepath.Join(dotfilePath, name)

  if DryRun {
    fmt.Printf("Will copy %s -> %s \n", configSrc, dotfileDest)
    return
  }

  _, err := fs.Stat(dotfileDest)
  if err == nil {
    fmt.Printf("Looks like %s exists in current dotfile directory\n", dotfileDest) 
    fmt.Println("Do you want to overwrite it?")
    confirm := promptui.Prompt{
      Label: "overwrite config",
      IsConfirm: true,
    }
    overwrite, _ := confirm.Run()
    if strings.ToUpper(overwrite) == "Y" {
      err = tools.CopyDir(fs, configSrc, dotfileDest)
      if err != nil {
        log.Fatal(err)
      }
      fmt.Printf("Copied %s -> %s\n", configSrc, dotfileDest)
    }
  }
}
