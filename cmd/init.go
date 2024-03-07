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

    var files []string
    var acceptedfiles [3] string 
    acceptedfiles[0] = "nvim"
    acceptedfiles[1] = "tmux"
    acceptedfiles[2] = "alacritty"

    err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        log.Fatalf("problem walking path %s", err)
        return nil
      }

      for _, acceptedfile := range(acceptedfiles) {
        if path == rootpath + acceptedfile {
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
      fmt.Fprintf(cmd.OutOrStdout(), file[len(rootpath):] + "\n" )
    }
    
  },
}
