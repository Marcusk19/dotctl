package tools

import (
	"fmt"
	"io"
	"os"
  "log"
	"path/filepath"
)

func CopyFile(srcFile, destFile string) error{
  // helper function to copy files over
  log.Printf("copy of %s to %s\n", srcFile, destFile)
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
  log.Printf("copying from %s to %s\n", srcDir, destDir)
  entries, err := os.ReadDir(srcDir)
  if err != nil {
    return err
  }
  log.Printf("entries found: %s\n", entries)

  for _, entry := range(entries) {
    if entry.Type().IsDir() {
      err := os.MkdirAll(filepath.Join(destDir, entry.Name()), os.ModePerm)
      if err != nil {
        return err
      }
      CopyDir(filepath.Join(srcDir, entry.Name()), filepath.Join(destDir, entry.Name()))
      continue
    }
    sourcePath := filepath.Join(srcDir, entry.Name())
    destPath := filepath.Join(destDir, entry.Name())

    err := CopyFile(sourcePath, destPath)
    if err != nil {
      return err
    }
  }

  return nil
}
