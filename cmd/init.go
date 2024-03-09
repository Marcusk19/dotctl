package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
  "github.com/Marcusk19/bender/tools"
)

func init() {
  RootCmd.AddCommand(initCommand)
}

func copyExistingConfigs(programs []string, destRootOpt ...string) {
  // takes list of programs and backs up configs for them
  destRoot := os.Getenv("HOME") + "/.dotfiles/"
  if len(destRootOpt) > 0 {
    destRoot = destRootOpt[0]
  }

  configRoot := os.Getenv("HOME") + "/.config/"
  for _, program := range(programs) {
    // TODO: do something here
    print(configRoot + program)
    err := tools.CopyDir(filepath.Join(configRoot, program), filepath.Join(destRoot, program))
    if err != nil {
      log.Fatal(err)
    }
  }
}

func createDotfileStructure(programs []string) {
  // takes list of programs and creates dotfiles for them
  dotfileRoot := os.Getenv("HOME") + "/.dotfiles/"
  for _, program := range(programs) {
    fmt.Printf("attempting to create directory %s%s\n", dotfileRoot, program)
    if err := os.MkdirAll(dotfileRoot + program, os.ModePerm); err != nil {
      log.Fatal(err)
    }
  }
}

var initCommand = &cobra.Command {
  Use: "init",
  Run: func(cmd *cobra.Command, args []string) {
    var rootpath string
    if len(args) <= 0 {
      fmt.Fprintf(cmd.OutOrStdout(), "no path provided, assuming /usr/bin/\n")
      rootpath = "/usr/bin/"
    } else {
      rootpath = args[0]
    }

    if rootpath[len(rootpath)-1:] != "/" {
      log.Fatal("path needs trailing slash")
    }

    var programs []string
    var acceptedprograms [3] string 
    acceptedprograms[0] = "nvim"
    acceptedprograms[1] = "tmux"
    acceptedprograms[2] = "alacritty"

    err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        log.Fatalf("problem walking path %s", err)
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

    fmt.Fprintf(cmd.OutOrStdout(), "binaries installed: \n =======================\n")
    for _, program := range(programs) {
      fmt.Fprintf(cmd.OutOrStdout(), program + "\n" )
    }

    createDotfileStructure(programs)
    copyExistingConfigs(programs)
    
  },
}
