package cmd

import (
	"fmt"
	"log"

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
    // we will do something here, for now just print args[0]
    fmt.Fprintf(cmd.OutOrStdout(), args[0])
  },
}
