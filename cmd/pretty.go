package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(prettyCmd)
}

var prettyCmd = &cobra.Command {
  Use: "pretty",
  Run: func(cmd *cobra.Command, args []string) {
    var filename = args[0]
    fmt.Printf("Run the pretty command with filename %s", filename)
  },
}
