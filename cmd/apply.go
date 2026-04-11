package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(applyCommand)
}

var applyCommand = &cobra.Command{
	Use:   "apply <repo-url>",
	Short: "Clone a dotfiles repo and link all tracked configs",
	Long: `apply clones a dotfiles repository (or pulls if it already exists),
reads the dotctl config, and creates symlinks for all tracked configs.

Example:
  dotctl apply https://github.com/user/dotfiles.git`,
	Args: cobra.ExactArgs(1),
	Run:  runApplyCommand,
}

func runApplyCommand(cmd *cobra.Command, args []string) {
	repoURL := args[0]
	dotfilePath := viper.GetString("dotfile-path")
	// strip trailing slash for consistency
	dotfilePath = filepath.Clean(dotfilePath)

	// Step 1: Clone or pull
	stat, err := os.Stat(dotfilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(cmd.OutOrStdout(), "Cloning %s into %s...\n", repoURL, dotfilePath)
		if !DryRun {
			_, err = git.PlainClone(dotfilePath, false, &git.CloneOptions{
				URL:      repoURL,
				Progress: cmd.OutOrStdout(),
			})
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: clone failed: %v\n", err)
				os.Exit(1)
			}
		}
	} else if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: cannot stat %s: %v\n", dotfilePath, err)
		os.Exit(1)
	} else if stat.IsDir() {
		// Check if it's a git repo
		repo, err := git.PlainOpen(dotfilePath)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s exists but is not a git repository\n", dotfilePath)
			fmt.Fprintf(cmd.ErrOrStderr(), "Remove it or use a different --dotfile-path\n")
			os.Exit(1)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Pulling latest changes in %s...\n", dotfilePath)
		if !DryRun {
			w, err := repo.Worktree()
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: cannot get worktree: %v\n", err)
				os.Exit(1)
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				fmt.Fprintf(cmd.OutOrStdout(), "Warning: git pull failed (%v), continuing with local state\n", err)
			} else if err == git.NoErrAlreadyUpToDate {
				fmt.Fprintln(cmd.OutOrStdout(), "Already up to date.")
			}
		}
	}

	// Step 2: Read config
	configPath := filepath.Join(dotfilePath, "dotctl", "config.yml")
	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: cannot read config at %s: %v\n", configPath, err)
		fmt.Fprintln(cmd.ErrOrStderr(), "Is this repo set up with dotctl? Run 'dotctl init' first.")
		os.Exit(1)
	}
	links := v.GetStringMapString("links")

	if len(links) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No links configured in config.yml — nothing to link.")
		return
	}

	// Step 3: Run idempotent link
	fmt.Fprintln(cmd.OutOrStdout(), "Linking dotfiles...")
	result := LinkDotfiles(cmd.OutOrStdout(), dotfilePath, links, Overwrite, NoBackup, DryRun)

	// Step 4: Print summary
	fmt.Fprintf(cmd.OutOrStdout(), "\nDone! %d linked, %d skipped, %d backed up\n",
		result.Linked, result.Skipped, result.Backed)
}
