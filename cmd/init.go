package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Marcusk19/bender/tools"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func init() {
  RootCmd.AddCommand(initCommand)
}

func copyExistingConfigs(programs []string, fs afero.Fs) {
  // takes list of programs and backs up configs for them
  destRoot := DotfilePath

  configRoot := ConfigPath
  for _, program := range(programs) {
    // TODO: do something here
    err := tools.CopyDir(fs, filepath.Join(configRoot, program), filepath.Join(destRoot, program))
    if err != nil {
      log.Fatalf("Problem copying %s", err.Error())
    }
  }
}

func createDotfileStructure(programs []string, fs afero.Fs) {
  // takes list of programs and creates dotfiles for them
  dotfileRoot := DotfilePath
  fmt.Printf("creating dotfile directory structure at %s\n", dotfileRoot)
  for _, program := range(programs) {
    if err := fs.MkdirAll(path.Join(dotfileRoot, program), os.ModePerm); err != nil {
      log.Fatal(err)
    }
  }
}

var initCommand = &cobra.Command {
  Use: "init",
  Short: "Copy configs to dotfile directory",
  Long: "Searches existing config directory for configs and then copies them to dotfile directory",
  Run: func(cmd *cobra.Command, args []string) {

    fs := FileSystem

    if(viper.Get("testing") == true && fs.Name() != "MemMapFS") {
      log.Fatalf("wrong filesystem, got %s", fs.Name())
    }

    var rootpath string
    if len(args) <= 0 {
      fmt.Fprintf(cmd.OutOrStdout(), "no path provided, assuming /usr/bin/\n")
      rootpath = "/usr/bin/"
    } else {
      rootpath = args[0]
    }

    if rootpath[len(rootpath)-1:] != "/" {
      log.Fatal("path needs trailing slash\n")
    }

    // TODO make a configurable list of binaries we want to look for
    var programs []string
    var acceptedprograms [3] string 
    acceptedprograms[0] = "nvim"
    acceptedprograms[1] = "tmux"
    acceptedprograms[2] = "alacritty"
    
    err := afero.Walk(fs, rootpath, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        log.Fatalf("problem walking path %s\n", err)
        return nil
      }

      for _, acceptedprogram := range(acceptedprograms) {
        if path == rootpath + acceptedprogram {
          programs = append(programs, path[len(rootpath):])
        }
      }
      return nil
    })

    if err != nil {
      log.Fatal(err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "binaries found: \n =======================\n")
    for _, program := range(programs) {
      fmt.Fprintf(cmd.OutOrStdout(), program + "\n" )
    }

    createDotfileStructure(programs, fs)
    copyExistingConfigs(programs, fs)

    if (viper.Get("testing") != true){
      _, err = git.PlainInit(DotfilePath, false)
      if err != nil {
        log.Fatal(err)
      }
    }
    fmt.Fprintf(cmd.OutOrStdout(), "Successfully created dotfiles repository\n")
  },
}
