package test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestInitCommand(t *testing.T) {
  viper.Set("testing", true)

  fs := cmd.FileSystem

  dotctl := cmd.RootCmd
  actual := new(bytes.Buffer)

  dotctl.SetOut(actual)
  dotctl.SetErr(actual)
  dotctl.SetArgs([]string{"init", "--dotfile-path=dotctl_test/dotfiles"})
  
  dotctl.Execute()

  homedir := "dotctl_test/"

  _, err := afero.ReadFile(fs, filepath.Join(homedir, "dotfiles/dotctl/config")) 
  if err != nil {
    t.Error(err.Error())
  }
  
}
