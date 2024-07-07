package test

import (
	"bytes"
	"os"
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
  dotctl.SetArgs([]string{"init"})
  
  dotctl.Execute()

  homedir := os.Getenv("HOME")

  _, err := afero.ReadFile(fs, filepath.Join(homedir, "dotfiles/dotctl/config.yml")) 
  if err != nil {
    t.Error(err.Error())
  }
  
}
