package test

import (
	"bytes"
	"testing"

	"github.com/Marcusk19/bender/cmd"
	"github.com/stretchr/testify/assert"
)

func TestPrettyCommand(t *testing.T) {
  bender := cmd.RootCmd
  actual := new(bytes.Buffer)
  bender.SetOut(actual)
  bender.SetErr(actual)
  bender.SetArgs([]string{"pretty", "fixtures/test_pretty.txt"})
  bender.Execute()
  
  expected := "The end of this sentence should start a newline. \nThe next sentence should be indented below this one.\n\tHello this is the end of the text"
  assert.Equal(t, expected, actual.String(), "actual value differs from expected")
}

