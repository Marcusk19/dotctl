package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LinkResult holds counts from a link operation
type LinkResult struct {
	Linked  int
	Skipped int
	Backed  int
}

func init() {
	RootCmd.AddCommand(linkCommand)
}

var linkCommand = &cobra.Command{
	Use:   "link",
	Run:   runLinkCommand,
	Short: "generate symlinks according to config",
	Long:  "runs through all configs in the dotctl config file and links them to configured symlinks",
}

func runLinkCommand(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "Symlinking dotfiles...")
	dotfileRoot := viper.GetString("dotfile-path")
	links := viper.GetStringMapString("links")

	result := LinkDotfiles(cmd.OutOrStdout(), dotfileRoot, links, Overwrite, NoBackup, DryRun)

	fmt.Fprintf(cmd.OutOrStdout(), "\nSummary: %d linked, %d skipped, %d backed up\n",
		result.Linked, result.Skipped, result.Backed)
}

// LinkDotfiles performs idempotent symlinking of dotfiles.
// dotfileRoot: path to the dotfiles repo (e.g. ~/dotfiles)
// links: map of config name -> target path (from config.yml)
// overwrite: if true, replaces existing wrong symlinks
// noBackup: if true, skips conflicts instead of backing up real files
// dryRun: if true, prints actions without making filesystem changes
func LinkDotfiles(out io.Writer, dotfileRoot string, links map[string]string, overwrite, noBackup, dryRun bool) LinkResult {
	result := LinkResult{}

	for configName, targetPath := range links {
		if configName == ".git" || configName == "dotctl" {
			continue
		}

		dotPath := filepath.Join(dotfileRoot, configName)

		if targetPath == "" {
			fmt.Fprintf(out, "  Warning: no target path configured for %s, skipping\n", configName)
			result.Skipped++
			continue
		}

		// Ensure parent directory exists
		parentDir := filepath.Dir(targetPath)
		if !dryRun {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				log.Printf("Warning: could not create parent dir %s: %v\n", parentDir, err)
			}
		}

		// Check what's at the target path using Lstat (doesn't follow symlinks)
		lstat, err := os.Lstat(targetPath)

		if os.IsNotExist(err) {
			// Target missing - create symlink
			if dryRun {
				fmt.Fprintf(out, "  %s → %s [would link]\n", targetPath, dotPath)
			} else {
				if err := os.Symlink(dotPath, targetPath); err != nil {
					log.Printf("Error: cannot symlink %s → %s: %v\n", targetPath, dotPath, err)
					result.Skipped++
					continue
				}
				fmt.Fprintf(out, "  %s → %s [linked]\n", targetPath, dotPath)
			}
			result.Linked++
			continue
		}

		if err != nil {
			log.Printf("Warning: cannot stat %s: %v\n", targetPath, err)
			result.Skipped++
			continue
		}

		// Check if it's a symlink
		if lstat.Mode()&os.ModeSymlink != 0 {
			existing, err := os.Readlink(targetPath)
			if err == nil && existing == dotPath {
				// Already pointing to the right place
				fmt.Fprintf(out, "  %s → %s [already linked]\n", targetPath, dotPath)
				result.Skipped++
				continue
			}
			// Symlink to wrong place
			if overwrite {
				if dryRun {
					fmt.Fprintf(out, "  %s [would remove wrong symlink and relink → %s]\n", targetPath, dotPath)
				} else {
					os.Remove(targetPath)
					if err := os.Symlink(dotPath, targetPath); err != nil {
						log.Printf("Error: cannot relink %s: %v\n", targetPath, err)
						result.Skipped++
						continue
					}
					fmt.Fprintf(out, "  %s → %s [relinked]\n", targetPath, dotPath)
				}
				result.Linked++
			} else {
				fmt.Fprintf(out, "  %s [symlink to wrong path, skipping (use --overwrite to force)]\n", targetPath)
				result.Skipped++
			}
			continue
		}

		// Real file or directory exists at target
		if noBackup {
			fmt.Fprintf(out, "  %s [real file exists, skipping (use --overwrite or remove --no-backup)]\n", targetPath)
			result.Skipped++
			continue
		}

		backupPath := targetPath + ".dotctl.bak"
		if dryRun {
			fmt.Fprintf(out, "  %s [would back up to %s and link → %s]\n", targetPath, backupPath, dotPath)
		} else {
			if err := os.Rename(targetPath, backupPath); err != nil {
				log.Printf("Error: cannot back up %s: %v\n", targetPath, err)
				result.Skipped++
				continue
			}
			if err := os.Symlink(dotPath, targetPath); err != nil {
				log.Printf("Error: cannot symlink after backup %s: %v\n", targetPath, err)
				result.Skipped++
				continue
			}
			fmt.Fprintf(out, "  %s → %s [backed up to .dotctl.bak, linked]\n", targetPath, dotPath)
		}
		result.Backed++
	}

	return result
}
