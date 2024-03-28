package test

import (
	"bytes"
	"testing"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/stretchr/testify/assert"
)

func TestPrettyCommand(t *testing.T) {
  dotctl := cmd.RootCmd
  actual := new(bytes.Buffer)
  dotctl.SetOut(actual)
  dotctl.SetErr(actual)
  dotctl.SetArgs([]string{"pretty", "fixtures/test_pretty.txt"})
  dotctl.Execute()
  
  expected := "The end of this sentence should start a newline. \nThe next sentence should be indented below this one.\n\tHello this is the end of the text"
  assert.Equal(t, expected, actual.String(), "actual value differs from expected")
}

