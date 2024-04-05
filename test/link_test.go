package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestLinkCommand(t *testing.T) {
  viper.Set("testing", true)
  cmd.FileSystem = afero.NewMemMapFs()
  fs := cmd.FileSystem
  homedir := os.Getenv("HOME")

  fs.MkdirAll(filepath.Join(homedir, "dotfiles/dotctl"), 0755)
  links := map[string]string {
    "someconfig": filepath.Join(homedir, ".config/someconfig"),
  }
  viper.Set("links", links)

  dotctl := cmd.RootCmd
  actual := new(bytes.Buffer)

  dotctl.SetOut(actual)
  dotctl.SetErr(actual)
  dotctl.SetArgs([]string{"link"})

  dotctl.Execute()

  someconfig := filepath.Join(homedir, ".config/someconfig/")
  somedot := filepath.Join(homedir, "dotfiles/someconfig/")

  expected := fmt.Sprintf("%s,%s", someconfig, somedot)

  assert.Equal(t, expected, actual.String(), "actual differs from expected")

}

