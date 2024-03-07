package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
  RootCmd.AddCommand(initCommand)
}

var initCommand = &cobra.Command {
  Use: "init",
  Run: func(cmd *cobra.Command, args []string) {
    if(len(args) <= 0) {
      log.Fatal(cmd.OutOrStdout(), "no arguments provided to init")
    }

    if args[0][len(args[0])-1:] != "/" {
      log.Fatal("path needs trailing slash")
    }

    var files []string
    var acceptedfiles [3] string 
    acceptedfiles[0] = "nvim"
    acceptedfiles[1] = "tmux"
    acceptedfiles[2] = "alacritty"

    rootpath := args[0]
    err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        log.Fatalf("problem walking path %s", err)
        return nil
      }

      for _, acceptedfile := range(acceptedfiles) {
        if path == args[0] + acceptedfile {
          files = append(files, path)
        }
      }
      return nil
    })

    if err != nil {
      log.Fatal(err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "binaries installed: \n =======================\n")
    for _, file := range(files) {
      fmt.Fprintf(cmd.OutOrStdout(), file[len(args[0]):] + "\n" )
    }
    
  },
}
