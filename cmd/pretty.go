package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
  RootCmd.AddCommand(prettyCmd)
}

var prettyCmd = &cobra.Command {
  Use: "pretty",
  Run: func(cmd *cobra.Command, args []string) {
    if (len(args) <= 0) {
      log.Fatal("no arguments provided")
    }
    var filename = args[0]
    f, err := os.Open(filename)
    if err != nil {
      log.Fatal(err)
    }

    defer f.Close()

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
      line := scanner.Text()
      formattedLine := strings.Replace(line, "\\n", "\n", -1)
      formattedLine = strings.Replace(formattedLine, "\\t", "\t", -1)
      fmt.Fprintf(cmd.OutOrStdout(), formattedLine)
    }

    if err := scanner.Err(); err != nil {
      log.Fatal(err)
    }
  },
}
