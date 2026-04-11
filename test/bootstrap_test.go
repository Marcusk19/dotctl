package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- DiscoverLinks unit tests ---

func TestDiscoverLinks_Basic(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "nvim"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "zsh"), 0755))

	configRoot := t.TempDir()
	links, err := cmd.DiscoverLinks(dir, configRoot)
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(configRoot, "nvim"), links["nvim"])
	assert.Equal(t, filepath.Join(configRoot, "zsh"), links["zsh"])
	assert.Len(t, links, 2)
}

func TestDiscoverLinks_SkipsGitAndDotctl(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".git"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "dotctl"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "nvim"), 0755))

	configRoot := t.TempDir()
	links, err := cmd.DiscoverLinks(dir, configRoot)
	require.NoError(t, err)

	assert.NotContains(t, links, ".git")
	assert.NotContains(t, links, "dotctl")
	assert.Contains(t, links, "nvim")
	assert.Len(t, links, 1)
}

func TestDiscoverLinks_Empty(t *testing.T) {
	dir := t.TempDir()
	configRoot := t.TempDir()

	links, err := cmd.DiscoverLinks(dir, configRoot)
	require.NoError(t, err)
	assert.Empty(t, links)
}

// --- WriteBootstrapConfig unit tests ---

func TestWriteBootstrapConfig_CreatesFile(t *testing.T) {
	dir := initLocalGitRepo(t)
	links := map[string]string{
		"nvim": "/home/user/.config/nvim",
		"zsh":  "/home/user/.config/zsh",
	}

	var buf bytes.Buffer
	err := cmd.WriteBootstrapConfig(&buf, dir, links, false)
	require.NoError(t, err)

	configFile := filepath.Join(dir, "dotctl", "config.yml")
	_, statErr := os.Stat(configFile)
	assert.NoError(t, statErr, "config.yml should exist")

	// Parse it back and verify
	v := viper.New()
	v.SetConfigFile(configFile)
	require.NoError(t, v.ReadInConfig())
	written := v.GetStringMapString("links")
	assert.Equal(t, "/home/user/.config/nvim", written["nvim"])
	assert.Equal(t, "/home/user/.config/zsh", written["zsh"])
}

func TestWriteBootstrapConfig_DryRun(t *testing.T) {
	dir := initLocalGitRepo(t)
	links := map[string]string{"nvim": "/home/user/.config/nvim"}

	var buf bytes.Buffer
	err := cmd.WriteBootstrapConfig(&buf, dir, links, true)
	require.NoError(t, err)

	configFile := filepath.Join(dir, "dotctl", "config.yml")
	_, statErr := os.Stat(configFile)
	assert.True(t, os.IsNotExist(statErr), "config.yml should NOT be created in dry-run mode")

	assert.Contains(t, buf.String(), "dry-run")
}

// --- bootstrap command integration tests ---

func TestBootstrapCommand_LinksCreated(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	configRoot := t.TempDir()

	// Put nvim and zsh dirs in the repo
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "zsh"), 0755))

	viper.Set("dotfile-path", repoDir)
	viper.Set("config-path", configRoot)
	defer viper.Set("dotfile-path", "")
	defer viper.Set("config-path", "")

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bootstrap", "https://fake.url/repo.git"})

	rootCmd.Execute()

	// Symlinks should exist at configRoot/nvim and configRoot/zsh
	nvimLink, err := os.Readlink(filepath.Join(configRoot, "nvim"))
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(repoDir, "nvim"), nvimLink)

	zshLink, err := os.Readlink(filepath.Join(configRoot, "zsh"))
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(repoDir, "zsh"), zshLink)
}

func TestBootstrapCommand_WritesConfigYml(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	configRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))

	viper.Set("dotfile-path", repoDir)
	viper.Set("config-path", configRoot)
	defer viper.Set("dotfile-path", "")
	defer viper.Set("config-path", "")

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bootstrap", "https://fake.url/repo.git"})

	rootCmd.Execute()

	configFile := filepath.Join(repoDir, "dotctl", "config.yml")
	_, err := os.Stat(configFile)
	assert.NoError(t, err, "config.yml should be written by bootstrap")

	v := viper.New()
	v.SetConfigFile(configFile)
	require.NoError(t, v.ReadInConfig())
	links := v.GetStringMapString("links")
	assert.Equal(t, filepath.Join(configRoot, "nvim"), links["nvim"])
}

func TestBootstrapCommand_DryRun(t *testing.T) {
	defer resetGlobalState()

	repoDir := initLocalGitRepo(t)
	configRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(repoDir, "nvim"), 0755))

	viper.Set("dotfile-path", repoDir)
	viper.Set("config-path", configRoot)
	defer viper.Set("dotfile-path", "")
	defer viper.Set("config-path", "")

	cmd.DryRun = true

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bootstrap", "https://fake.url/repo.git"})

	rootCmd.Execute()

	output := buf.String()
	assert.Contains(t, output, "would link")
	assert.Contains(t, output, "dry-run")

	// No symlink should be created
	_, err := os.Lstat(filepath.Join(configRoot, "nvim"))
	assert.True(t, os.IsNotExist(err), "symlink should NOT be created in dry-run mode")

	// No config.yml should be written
	_, err = os.Stat(filepath.Join(repoDir, "dotctl", "config.yml"))
	assert.True(t, os.IsNotExist(err), "config.yml should NOT be written in dry-run mode")
}

func TestBootstrapCommand_NoArgs(t *testing.T) {
	defer resetGlobalState()

	rootCmd := cmd.RootCmd
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bootstrap"})

	err := rootCmd.Execute()
	assert.Error(t, err, "bootstrap with no args should return an error")
}
