package test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/bender/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestInitCommand(t *testing.T) {
  viper.Set("testing", true)

  fs := cmd.FileSystem

  bender := cmd.RootCmd
  actual := new(bytes.Buffer)

  bender.SetOut(actual)
  bender.SetErr(actual)
  bender.SetArgs([]string{"init", "--dotfile-path=bender_test/dotfiles"})
  
  bender.Execute()

  homedir := "bender_test/"

  _, err := afero.ReadFile(fs, filepath.Join(homedir, "dotfiles/bender/config")) 
  if err != nil {
    t.Error(err.Error())
  }
  
}
