package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/bender/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestLinkCommand(t *testing.T) {
  setUpTesting()
  bender := cmd.RootCmd
  actual := new(bytes.Buffer)

  bender.SetOut(actual)
  bender.SetErr(actual)
  bender.SetArgs([]string{"link"})

  bender.Execute()

  homedir := os.Getenv("HOME")
  someconfig := filepath.Join(homedir, ".config/someconfig/")
  somedot := filepath.Join(homedir, ".dotfiles/someconfig/")

  expected := fmt.Sprintf("%s,%s", someconfig, somedot)

  assert.Equal(t, expected, actual.String(), "actual differs from expected")

  tearDownTesting()
}

func setUpTesting() {
  fs := cmd.FileSystem
  homedir := os.Getenv("HOME")
  fs.MkdirAll(filepath.Join(homedir, ".dotfiles/bender"), 0755)
  fs.Create(filepath.Join(homedir, ".dotfiles/bender/config"))
  fs.MkdirAll(filepath.Join(homedir, ".dotfiles/someconfig/"), 0755)

  viper.Set("dotfile-path", filepath.Join(homedir, ".dotfiles"))
  viper.Set("someconfig", filepath.Join(homedir, ".config/someconfig/"))
  viper.Set("testing", true)

}

func tearDownTesting() {
  fs := cmd.FileSystem
  fs.RemoveAll("bender_test/")
}
