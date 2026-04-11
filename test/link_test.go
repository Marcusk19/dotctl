package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkDotfiles_MissingTarget(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	// Create source dir in dotfileRoot
	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	// Target does not exist yet
	targetPath := filepath.Join(targetRoot, ".config", "nvim")

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 1, result.Linked, "should have linked 1")
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 0, result.Backed)

	// Verify symlink was actually created and points to the right place
	linkTarget, err := os.Readlink(targetPath)
	require.NoError(t, err, "symlink should exist at targetPath")
	assert.Equal(t, dotPath, linkTarget, "symlink should point to dotPath")
}

func TestLinkDotfiles_AlreadyLinked(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	require.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0755))

	// Pre-create correct symlink
	require.NoError(t, os.Symlink(dotPath, targetPath))

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 1, result.Skipped, "should have skipped 1 (already linked)")
	assert.Equal(t, 0, result.Backed)

	// Symlink should still point to dotPath
	linkTarget, err := os.Readlink(targetPath)
	require.NoError(t, err)
	assert.Equal(t, dotPath, linkTarget)

	assert.Contains(t, buf.String(), "already linked")
}

func TestLinkDotfiles_WrongSymlink_NoOverwrite(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	wrongPath := filepath.Join(t.TempDir(), "wrong")
	require.NoError(t, os.MkdirAll(wrongPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	require.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0755))

	// Pre-create symlink pointing to the wrong place
	require.NoError(t, os.Symlink(wrongPath, targetPath))

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 1, result.Skipped, "should have skipped 1 (wrong symlink, no overwrite)")
	assert.Equal(t, 0, result.Backed)

	// Symlink should still point to the wrong path (unchanged)
	linkTarget, err := os.Readlink(targetPath)
	require.NoError(t, err)
	assert.Equal(t, wrongPath, linkTarget)
}

func TestLinkDotfiles_WrongSymlink_Overwrite(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	wrongPath := filepath.Join(t.TempDir(), "wrong")
	require.NoError(t, os.MkdirAll(wrongPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	require.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0755))

	// Pre-create symlink pointing to the wrong place
	require.NoError(t, os.Symlink(wrongPath, targetPath))

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, true, false, false)

	assert.Equal(t, 1, result.Linked, "should have relinked 1")
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 0, result.Backed)

	// Symlink should now point to dotPath
	linkTarget, err := os.Readlink(targetPath)
	require.NoError(t, err)
	assert.Equal(t, dotPath, linkTarget, "symlink should now point to dotPath")
}

func TestLinkDotfiles_RealFile_WithBackup(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	require.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0755))

	// Create a real file at the target path
	require.NoError(t, os.WriteFile(targetPath, []byte("original content"), 0644))

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 1, result.Backed, "should have backed up 1")

	// Backup file should exist with original content
	backupPath := targetPath + ".dotctl.bak"
	backupContent, err := os.ReadFile(backupPath)
	require.NoError(t, err, "backup file should exist")
	assert.Equal(t, "original content", string(backupContent))

	// Target should now be a symlink to dotPath
	linkTarget, err := os.Readlink(targetPath)
	require.NoError(t, err, "target should be a symlink now")
	assert.Equal(t, dotPath, linkTarget)
}

func TestLinkDotfiles_RealFile_NoBackup(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	require.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0755))

	// Create a real file at the target path
	require.NoError(t, os.WriteFile(targetPath, []byte("original content"), 0644))

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, true, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 1, result.Skipped, "should have skipped 1 (noBackup)")
	assert.Equal(t, 0, result.Backed)

	// Original file should be untouched
	content, err := os.ReadFile(targetPath)
	require.NoError(t, err)
	assert.Equal(t, "original content", string(content))

	// No backup file should exist
	backupPath := targetPath + ".dotctl.bak"
	_, err = os.Stat(backupPath)
	assert.True(t, os.IsNotExist(err), "backup file should not exist")
}

func TestLinkDotfiles_DryRun(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	dotPath := filepath.Join(dotfileRoot, "nvim")
	require.NoError(t, os.MkdirAll(dotPath, 0755))

	targetPath := filepath.Join(targetRoot, ".config", "nvim")
	// Do NOT create parent dirs - dryRun should not create them either

	links := map[string]string{"nvim": targetPath}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, true)

	assert.Equal(t, 1, result.Linked, "should count 1 linked even in dry run")
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 0, result.Backed)

	// No symlink should actually be created
	_, err := os.Lstat(targetPath)
	assert.True(t, os.IsNotExist(err), "symlink should NOT be created in dry run")

	assert.Contains(t, buf.String(), "would link")
}

func TestLinkDotfiles_SkipsGitAndDotctl(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	links := map[string]string{
		".git":   filepath.Join(targetRoot, ".git"),
		"dotctl": filepath.Join(targetRoot, "dotctl"),
	}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 0, result.Backed)

	// Nothing should be created
	_, err := os.Lstat(filepath.Join(targetRoot, ".git"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Lstat(filepath.Join(targetRoot, "dotctl"))
	assert.True(t, os.IsNotExist(err))
}

func TestLinkDotfiles_EmptyTargetPath(t *testing.T) {
	dotfileRoot := t.TempDir()

	links := map[string]string{"nvim": ""}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 0, result.Linked)
	assert.Equal(t, 1, result.Skipped, "empty target should be skipped")
	assert.Equal(t, 0, result.Backed)
}

func TestLinkDotfiles_MultipleEntries(t *testing.T) {
	dotfileRoot := t.TempDir()
	targetRoot := t.TempDir()

	// Create two source dirs
	require.NoError(t, os.MkdirAll(filepath.Join(dotfileRoot, "nvim"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dotfileRoot, "zsh"), 0755))

	nvimTarget := filepath.Join(targetRoot, ".config", "nvim")
	zshTarget := filepath.Join(targetRoot, ".zshrc")

	links := map[string]string{
		"nvim": nvimTarget,
		"zsh":  zshTarget,
	}
	var buf bytes.Buffer
	result := cmd.LinkDotfiles(&buf, dotfileRoot, links, false, false, false)

	assert.Equal(t, 2, result.Linked, "should have linked 2 entries")
	assert.Equal(t, 0, result.Skipped)
	assert.Equal(t, 0, result.Backed)

	// Verify both symlinks
	linkTarget, err := os.Readlink(nvimTarget)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dotfileRoot, "nvim"), linkTarget)

	linkTarget, err = os.Readlink(zshTarget)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dotfileRoot, "zsh"), linkTarget)
}
