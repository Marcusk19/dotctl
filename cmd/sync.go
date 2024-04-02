package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/manifoldco/promptui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var remoteRepository string

func init() {
	RootCmd.AddCommand(syncCommand)
  syncCommand.Flags().StringVarP(
    &remoteRepository,
    "remote",
    "r",
    "",
    "URL of remote repository",
  )

  viper.BindPFlag("dotctl-origin", syncCommand.Flags().Lookup("remote"))
}

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotfiles with git",
	Long:  "TODO: add longer description",
	Run:   runSyncCommand,
}

func validateInput(input string) error {
  if input == ""  {
    return errors.New("Missing input")
  }

  return nil
}

func gitAddFiles(worktree *git.Worktree, fs afero.Fs) error {
  dotfilepath := viper.GetString("dotfile-path")
  entries, err := afero.ReadDir(fs, dotfilepath)
  if err != nil {
    return err
  }
  for _, entry := range(entries) {
    if(entry.Name() == "dotctl") {
      continue
    }
    _, err = worktree.Add(entry.Name())
    if err != nil {
      return err
    }
  }
  return nil
}

func runSyncCommand(cmd *cobra.Command, args []string) {
	origin := viper.GetString("dotctl-origin")
	if origin == "" {
		fmt.Fprintln(cmd.OutOrStdout(), "No remote repository found")
		return
	}

	dotfilepath := viper.GetString("dotfile-path")
	r, err := git.PlainOpen(dotfilepath)
	CheckIfError(err)

	// check remotes and if origin does not exist
	// we need to create it
	list, err := r.Remotes()
	CheckIfError(err)

	if len(list) == 0 {
		r.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{origin},
		})
	}

	w, err := r.Worktree()
	CheckIfError(err)

  username := promptui.Prompt{
    Label: "username",
    Validate: validateInput,
  }

  password := promptui.Prompt{
    Label: "password",
    Validate: validateInput,
  }

  usernameVal, err := username.Run()
  CheckIfError(err)

  passwordVal, err := password.Run()
  CheckIfError(err)

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
    Auth: &http.BasicAuth {
      Username: usernameVal,
      Password: passwordVal,
    },
	})

  if err != nil{
    fmt.Println(err)
  } else {
    fmt.Fprintf(cmd.OutOrStdout(), "successfully pulled from %s", origin)
  }

  err = gitAddFiles(w, FileSystem)
  if err != nil {
    log.Fatalf("Could not add files: %s\n", err)
  }
  
  commitMessage := "backup " + time.Now().String()

  commit, err := w.Commit(commitMessage, &git.CommitOptions{
    Author: &object.Signature{
      Name: "dotctl CLI",
      Email: "example@example.com",
      When: time.Now(),
    },
  })

  if err != nil {
    log.Fatal(err.Error())
  }

  obj, err := r.CommitObject(commit)

  if err != nil {
    log.Fatalf("Cannot commit: %s",err)
  }

  fmt.Println(obj)

  err = r.Push(&git.PushOptions{
		RemoteName: "origin",
    Auth: &http.BasicAuth {
      Username: usernameVal,
      Password: passwordVal,
    },
  })
  CheckIfError(err)
  viper.WriteConfig()

}
