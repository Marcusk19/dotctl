/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"time"

	"github.com/Marcusk19/dotctl/cmd"
	"github.com/carlmjohnson/versioninfo"
)

var (
  version = "dev"
  commit = "none"
  date = "unknown"
)

func SetVersionInfo(version, commit, date string) {
  cmd.RootCmd.Version = fmt.Sprintf("%s [Built on %s from Git Sha %s]", version, date, commit)
}

func main() {
  SetVersionInfo(versioninfo.Version, versioninfo.Revision, versioninfo.LastCommit.Format(time.RFC3339))
	cmd.Execute()
}
