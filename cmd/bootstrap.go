package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(bootstrapCommand)
}

var bootstrapCommand = &cobra.Command{
	Use:   "bootstrap <repo-url>",
	Short: "Clone a dotfiles repo and auto-discover configs to link (no config.yml required)",
	Long: `bootstrap clones a dotfiles repository (or pulls if it already exists),
auto-discovers all top-level directories and files (excluding .git and dotctl),
symlinks each one to ~/.config/<name>, and writes a config.yml so that
subsequent dotctl commands (status, rm, link) work normally.

Use this when your dotfiles repo mirrors ~/.config/ directly and does not
have a dotctl config.yml yet.

Example:
  dotctl bootstrap https://github.com/user/dotfiles.git`,
	Args: cobra.ExactArgs(1),
	Run:  runBootstrapCommand,
}

func runBootstrapCommand(cmd *cobra.Command, args []string) {
	repoURL := args[0]
	dotfilePath := filepath.Clean(viper.GetString("dotfile-path"))
	configPath := filepath.Clean(viper.GetString("config-path"))

	// Step 1: Clone or pull
	if err := cloneOrPull(cmd, repoURL, dotfilePath, DryRun); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Discover links from repo structure
	links, err := DiscoverLinks(dotfilePath, configPath)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: cannot discover links: %v\n", err)
		os.Exit(1)
	}

	if len(links) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No configs discovered in repo — nothing to link.")
		return
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Discovered %d config(s) to link:\n", len(links))
	for name, target := range links {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s → %s\n", name, target)
	}
	fmt.Fprintln(cmd.OutOrStdout())

	// Step 3: Link dotfiles
	fmt.Fprintln(cmd.OutOrStdout(), "Linking dotfiles...")
	result := LinkDotfiles(cmd.OutOrStdout(), dotfilePath, links, Overwrite, NoBackup, DryRun)

	// Step 4: Write config.yml so future dotctl commands work
	if err := WriteBootstrapConfig(cmd.OutOrStdout(), dotfilePath, links, DryRun); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Warning: could not write config.yml: %v\n", err)
	}

	// Step 5: Summary
	fmt.Fprintf(cmd.OutOrStdout(), "\nDone! %d linked, %d skipped, %d backed up\n",
		result.Linked, result.Skipped, result.Backed)
}

// DiscoverLinks scans dotfileRoot for top-level entries, skipping .git and dotctl,
// and returns a links map of name → filepath.Join(configRoot, name).
func DiscoverLinks(dotfileRoot, configRoot string) (map[string]string, error) {
	entries, err := os.ReadDir(dotfileRoot)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %w", dotfileRoot, err)
	}

	links := make(map[string]string)
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" || name == "dotctl" {
			continue
		}
		links[name] = filepath.Join(configRoot, name)
	}
	return links, nil
}

// WriteBootstrapConfig writes a config.yml containing the discovered links
// to dotfileRoot/dotctl/config.yml using a fresh viper instance.
// If dryRun is true, it only prints what would be written.
func WriteBootstrapConfig(out io.Writer, dotfileRoot string, links map[string]string, dryRun bool) error {
	configDir := filepath.Join(dotfileRoot, "dotctl")
	configFile := filepath.Join(configDir, "config.yml")

	if dryRun {
		fmt.Fprintf(out, "[dry-run] Would write config.yml with %d entries to %s\n", len(links), configFile)
		return nil
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("cannot create dotctl dir: %w", err)
	}

	v := viper.New()
	v.Set("links", links)
	if err := v.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}

	fmt.Fprintf(out, "\nWrote config.yml to %s\n", configFile)
	return nil
}
