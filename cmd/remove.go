package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Marcusk19/dotctl/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(removeCommand)
}

var removeCommand = &cobra.Command{
	Use:   "rm",
	Short: "remove dotfile link",
	Long:  "TODO: add longer description",
	Run:   runRemoveCommand,
}

func runRemoveCommand(cmd *cobra.Command, args []string) {
	fs := FileSystem

	if len(args) <= 0 {
		fmt.Println("ERROR: missing specified config")
		return
	}

	dotfile := args[0]
	links := viper.GetStringMapString("links")
	dotfileConfigPath := links[dotfile]

	err := fs.Remove(dotfileConfigPath)
	if err != nil {
		fmt.Printf("ERROR: problem removing symlink %s: %s\n", dotfileConfigPath, err)
		return
	}

	dotfileSavedPath := filepath.Join(DotfilePath, dotfile)
	savedFile, err := fs.Open(dotfileSavedPath)
	if err != nil {
		fmt.Printf("ERROR: problem viewing saved dotfile(s): %s\n", err)
		return
	}

	fileInfo, err := savedFile.Stat()
	if err != nil {
		fmt.Printf("ERROR: problem getting file info: %s\n", err)
		return
	}
	if fileInfo.IsDir() {
		err = tools.CopyDir(fs, dotfileSavedPath, dotfileConfigPath)
	} else {
		err = tools.CopyFile(fs, dotfileSavedPath, dotfileConfigPath)
	}

	if err != nil {
		fmt.Printf("ERROR: problem copying over dotfile(s) %s\n", err)
		return
	}

	delete(links, dotfile)
	viper.Set("links", links)
	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("ERROR: problem saving config: %s\n", err)
		return
	}

	fmt.Printf("%s symlink removed, copied files over to %s\n", dotfile, dotfileConfigPath)
}
