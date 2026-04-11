package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// initLocalGitRepo creates a temp directory and initializes it as a git repo.
func initLocalGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	_, err := gogit.PlainInit(dir, false)
	require.NoError(t, err)
	return dir
}

// writeApplyConfig writes a dotctl/config.yml with the given links map.
func writeApplyConfig(t *testing.T, dir string, links map[string]string) {
	t.Helper()
	configDir := filepath.Join(dir, "dotctl")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	content := "links:\n"
	for name, path := range links {
		content += fmt.Sprintf("  %s: %s\n", name, path)
	}
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(content), 0644))
}

// resetGlobalState resets the global cmd flags to defaults.
func resetGlobalState() {
	cmd.DryRun = false
	cmd.Overwrite = false
	cmd.NoBackup = false
}

func TestApplyCommand_ExistingRepo_EmptyLinks(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	writeApplyConfig(t, repoDir, map[string]string{})

	viper.Set("dotfile-path", repoDir)
	defer viper.Set("dotfile-path", "")

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"apply", "https://fake.url/repo.git"})

	rootCmd.Execute()

	output := buf.String()
	assert.Contains(t, output, "Pulling latest changes")
	assert.Contains(t, output, "No links configured")
}

func TestApplyCommand_ExistingRepo_WithLinks_DryRun(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	targetRoot := t.TempDir()

	nvimTarget := filepath.Join(targetRoot, ".config", "nvim")

	writeApplyConfig(t, repoDir, map[string]string{
		"nvim": nvimTarget,
	})

	// Create source directory in the repo
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))

	viper.Set("dotfile-path", repoDir)
	defer viper.Set("dotfile-path", "")

	cmd.DryRun = true

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"apply", "https://fake.url/repo.git"})

	rootCmd.Execute()

	output := buf.String()
	assert.Contains(t, output, "Linking dotfiles...")
	assert.Contains(t, output, "would link")

	// Verify no actual symlink was created
	_, err := os.Lstat(nvimTarget)
	assert.True(t, os.IsNotExist(err), "symlink should NOT be created in dry-run mode")
}

func TestApplyCommand_ExistingRepo_WithLinks(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	targetRoot := t.TempDir()

	nvimTarget := filepath.Join(targetRoot, ".config", "nvim")

	writeApplyConfig(t, repoDir, map[string]string{
		"nvim": nvimTarget,
	})

	// Create source directory in the repo
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))

	viper.Set("dotfile-path", repoDir)
	defer viper.Set("dotfile-path", "")

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"apply", "https://fake.url/repo.git"})

	rootCmd.Execute()

	output := buf.String()
	assert.Contains(t, output, "Linking dotfiles...")
	assert.Contains(t, output, "1 linked, 0 skipped, 0 backed up")

	// Verify symlink was actually created
	linkTarget, err := os.Readlink(nvimTarget)
	require.NoError(t, err, "symlink should exist at target")
	assert.Equal(t, filepath.Join(repoDir, "nvim"), linkTarget)
}

func TestApplyCommand_NoArgs(t *testing.T) {
	defer resetGlobalState()

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"apply"})

	err := rootCmd.Execute()

	assert.Error(t, err, "apply with no args should return an error")
}

func TestApplyCommand_ExistingRepo_MultipleLinks_DryRun(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	targetRoot := t.TempDir()

	nvimTarget := filepath.Join(targetRoot, ".config", "nvim")
	zshTarget := filepath.Join(targetRoot, ".zshrc")

	writeApplyConfig(t, repoDir, map[string]string{
		"nvim": nvimTarget,
		"zsh":  zshTarget,
	})

	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "zsh"), 0755))

	viper.Set("dotfile-path", repoDir)
	defer viper.Set("dotfile-path", "")

	cmd.DryRun = true

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"apply", "https://fake.url/repo.git"})

	rootCmd.Execute()

	output := buf.String()
	assert.Contains(t, output, "Linking dotfiles...")
	assert.Contains(t, output, "would link")
	assert.Contains(t, output, "Done!")
}
