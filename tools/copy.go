package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(srcFile, destFile string) error{
  // helper function to copy files over
  sourceFileStat, err := os.Stat(srcFile)
  if err != nil {
    return err
  }

  if !sourceFileStat.Mode().IsRegular() {
    return fmt.Errorf("%s is not a regular file", srcFile)
  }

  source, err := os.Open(srcFile)
  if err != nil {
    return err
  }
  defer source.Close()

  destination, err := os.Create(destFile)
  if err != nil {
    return err
  }
  defer destination.Close()
  _, err = io.Copy(destination, source)

  return err
}

func CopyDir(srcDir, destDir string) error {
  entries, err := os.ReadDir(srcDir)
  if err != nil {
    return err
  }

  for _, entry := range(entries) {
    sourcePath := filepath.Join(srcDir, entry.Name())
    destPath := filepath.Join(destDir, entry.Name())

    err := CopyFile(sourcePath, destPath)
    if err != nil {
      return err
    }
  }

  return nil
}
