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
  print("setting test var\n")
  viper.Set("testing", true)

  fs := cmd.SetUpForTesting()

  bender := cmd.RootCmd
  actual := new(bytes.Buffer)

  bender.SetOut(actual)
  bender.SetErr(actual)
  bender.SetArgs([]string{"init", "bin/", "--dotfile-path=bender_test/.dotfiles", "--config-path=bender_test/.config"})
  
  bender.Execute()

  homedir := "bender_test/"

  _, err := afero.ReadFile(fs, filepath.Join(homedir, ".dotfiles/alacritty/alacritty.conf")) 
  if err != nil {
    t.Error(err.Error())
  }
  
}
