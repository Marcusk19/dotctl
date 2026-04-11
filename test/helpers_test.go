package test

import (
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	gogit "github.com/go-git/go-git/v5"
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

// resetGlobalState resets the global cmd flags to defaults.
func resetGlobalState() {
	cmd.DryRun = false
	cmd.Overwrite = false
	cmd.NoBackup = false
}
