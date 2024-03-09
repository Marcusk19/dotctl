package tools

import (
	"log"

	"github.com/spf13/afero"
)

var AppFs afero.Fs = afero.NewOsFs()


func SetTestFs() {
  log.Print("setting test fs")
  testFs := afero.NewMemMapFs()
  AppFs = testFs
}
