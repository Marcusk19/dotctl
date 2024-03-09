package tools

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

func CopyFile(os afero.Fs, srcFile, destFile string) error{
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

func CopyDir(os afero.Fs, srcDir, destDir string) error {
  entries, err := afero.ReadDir(os, srcDir)
  if err != nil {
    return err
  }

  for _, entry := range(entries) {
    if entry.IsDir() {
      subDir := filepath.Join(srcDir, entry.Name())
      destSubDir := filepath.Join(destDir, entry.Name())
      err := os.MkdirAll(destSubDir, entry.Mode().Perm())
      if err != nil {
        return err
      }
      CopyDir(os, subDir, destSubDir)
      continue
    }
    sourcePath := filepath.Join(srcDir, entry.Name())
    destPath := filepath.Join(destDir, entry.Name())

    err := CopyFile(os, sourcePath, destPath)
    if err != nil {
      return err
    }
  }

  return nil
}