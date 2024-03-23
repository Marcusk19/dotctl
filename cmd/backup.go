package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Marcusk19/bender/tools"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func init() {
  RootCmd.AddCommand(backupCommand)
}

var backupCommand = &cobra.Command {
  Use: "backup",
  Short: "Add and commit files in dotfiles directory",
  Run: runBackup,
}


func runBackup(cmd *cobra.Command, args []string) {
  fmt.Fprintf(cmd.OutOrStdout(), "Backing up %s...\n", DotfilePath)
  r, err := git.PlainOpen(DotfilePath)
  if err != nil {
    log.Fatal(err)
  }
  
  worktree, err := r.Worktree()
  if err != nil {
    log.Fatal(err)
  }
  gitAddFiles(worktree, tools.AppFs)
  
  commitMessage := "backup " + time.Now().String()

  commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
    Author: &object.Signature{
      Name: "bender CLI",
      Email: "example@example.com",
      When: time.Now(),
    },
  })

  if err != nil {
    log.Fatal(err.Error())
  }

  obj, err := r.CommitObject(commit)

  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(obj)
}
