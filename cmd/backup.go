package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Marcusk19/bender/tools"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
  RootCmd.AddCommand(backupCommand)
  RootCmd.PersistentFlags().StringVar(&DotfilePath, "dotfile-path", "~/.dotfiles", "Path pointing to dotfiles directory")
  RootCmd.MarkFlagRequired("dotfile-path")
  viper.BindPFlag("dotfile-path", RootCmd.PersistentFlags().Lookup("dotfile-path"))
}

var backupCommand = &cobra.Command {
  Use: "backup",
  Run: runBackup,
}

var DotfilePath string

func gitAddFiles(worktree *git.Worktree, fs afero.Fs) error {
  entries, err := afero.ReadDir(fs, DotfilePath)
  if err != nil {
    return err
  }
  for _, entry := range(entries) {
    _, err = worktree.Add(entry.Name())
    if err != nil {
      return err
    }
  }
  return nil
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
