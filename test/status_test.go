package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStatusCommand(t *testing.T) {
  viper.Set("testing", true)

  fs := cmd.FileSystem

  homedir := os.Getenv("HOME")
  fs.MkdirAll(filepath.Join(homedir, "dotfiles/dotctl"), 0755)
  fs.MkdirAll(filepath.Join(homedir, "dotfiles/someconfig"), 0755)
  fs.MkdirAll(filepath.Join(homedir, "dotfiles/somelinkedconfig"), 0755)

  var links = map[string]string {
    "somelinkedconfig": "configpath",
  }
  
  viper.Set("links", links)

  dotctl := cmd.RootCmd

  actual := new(bytes.Buffer)

  dotctl.SetOut(actual)
  dotctl.SetErr(actual)
  dotctl.SetArgs([]string{"status"})

  dotctl.Execute()

  expected := "Config directories currently in dotfile path:\n" +
              "someconfig\nsomelinkedconfig - configpath\n"

  assert.Equal(t, expected, actual.String(), "actual differs from expected")
}
