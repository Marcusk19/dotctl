package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// cloneOrPull clones repoURL into dotfilePath if it does not exist,
// or pulls the latest changes if it is already a git repo.
// Does nothing on disk when dryRun is true.
func cloneOrPull(cmd *cobra.Command, repoURL, dotfilePath string, dryRun bool) error {
	stat, err := os.Stat(dotfilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(cmd.OutOrStdout(), "Cloning %s into %s...\n", repoURL, dotfilePath)
		if !dryRun {
			_, err = git.PlainClone(dotfilePath, false, &git.CloneOptions{
				URL:      repoURL,
				Progress: cmd.OutOrStdout(),
			})
			if err != nil {
				return fmt.Errorf("clone failed: %w", err)
			}
		}
	} else if err != nil {
		return fmt.Errorf("cannot stat %s: %w", dotfilePath, err)
	} else if stat.IsDir() {
		repo, err := git.PlainOpen(dotfilePath)
		if err != nil {
			return fmt.Errorf("%s exists but is not a git repository\nRemove it or use a different --dotfile-path", dotfilePath)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Pulling latest changes in %s...\n", dotfilePath)
		if !dryRun {
			w, err := repo.Worktree()
			if err != nil {
				return fmt.Errorf("cannot get worktree: %w", err)
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				fmt.Fprintf(cmd.OutOrStdout(), "Warning: git pull failed (%v), continuing with local state\n", err)
			} else if err == git.NoErrAlreadyUpToDate {
				fmt.Fprintln(cmd.OutOrStdout(), "Already up to date.")
			}
		}
	}
	return nil
}
