package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestLinkCommand(t *testing.T) {
  oldDotfilePath := viper.GetString("dotfile-path")
  setUpTesting()
  dotctl := cmd.RootCmd
  actual := new(bytes.Buffer)

  dotctl.SetOut(actual)
  dotctl.SetErr(actual)
  dotctl.SetArgs([]string{"link"})

  dotctl.Execute()

  homedir := os.Getenv("HOME")
  someconfig := filepath.Join(homedir, ".config/someconfig/")
  somedot := filepath.Join(homedir, "dotfiles/someconfig/")

  expected := fmt.Sprintf("%s,%s", someconfig, somedot)

  assert.Equal(t, expected, actual.String(), "actual differs from expected")

  tearDownTesting(oldDotfilePath)
}

func setUpTesting() {
  viper.Set("testing", true)

  fs := cmd.FileSystem
  homedir := os.Getenv("HOME")
  fakeLinks := map[string]string {"someconfig": filepath.Join(homedir, ".config/someconfig")}
  viper.Set("links", fakeLinks)
  fs.MkdirAll(filepath.Join(homedir, "dotfiles/dotctl"), 0755)
  fs.Create(filepath.Join(homedir, "dotfiles/dotctl/config"))

  viper.Set("dotfile-path", filepath.Join(homedir, "dotfiles"))
  viper.Set("someconfig", filepath.Join(homedir, ".config/someconfig/"))
}

func tearDownTesting(oldDotfilePath string) {
  viper.Set("dotfile-path", oldDotfilePath)
  viper.WriteConfig()
}
