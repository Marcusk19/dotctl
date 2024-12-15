package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Marcusk19/dotctl/tools"
	"github.com/manifoldco/promptui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	addCommand.Flags().BoolVar(&absolutePath, "absolute", false, "absolute path of config")
	RootCmd.AddCommand(addCommand)
}

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Adds config to be tracked by dotctl",
	Long:  "will copy files passed as argument to the dotfiles directory and symlink them", // TODO add more description
	Run:   runAddCommand,
}

var absolutePath bool

func runAddCommand(cmd *cobra.Command, args []string) {
	fs := FileSystem

	testing := viper.GetBool("testing")

	if len(args) <= 0 {
		fmt.Println("ERROR: requires config path")
		return
	}

	configSrc := args[0]

	if !absolutePath {
		cwd, _ := os.Getwd()
		configSrc = cwd + "/" + configSrc
	}

	dirs := strings.Split(configSrc, "/")
	name := dirs[len(dirs)-1] // take the last section of the path, this should be the name
	if name[0] == '.' {
		name = name[1:]
	}

	links := viper.GetStringMap("links")
	links[name] = configSrc
	viper.Set("links", links)

	dotfilePath := viper.Get("dotfile-path").(string)

	dotfileDest := filepath.Join(dotfilePath, name)

	if DryRun {
		fmt.Printf("Will copy %s -> %s \n", configSrc, dotfileDest)
		return
	}

	_, err := fs.Stat(dotfileDest)
	if err == nil {
		fmt.Printf("Looks like %s exists in current dotfile directory\n", dotfileDest)
		fmt.Printf("Do you want to overwrite it with what is in %s?\n", configSrc)
		confirm := promptui.Prompt{
			Label:     "overwrite config",
			IsConfirm: true,
		}
		overwrite, _ := confirm.Run()
		if strings.ToUpper(overwrite) == "Y" {
			addConfigToDir(fs, configSrc, dotfileDest)
		}
	} else {
		addConfigToDir(fs, configSrc, dotfileDest)
	}

	if !DryRun {
		// symlink the copied dotfile destination back to the config src
		fs.RemoveAll(configSrc)
		linkPaths(dotfileDest, configSrc)
	} else {
		fmt.Println("Files were not symlinked")
	}

	if !testing {
		// write to the config to persist changes
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("Problem updating dotctl config %s", err)
		}
	}
}

func addConfigToDir(fs afero.Fs, configSrc, dotfileDest string) {
	configFile, err := fs.Open(configSrc)
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	fileInfo, err := configFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		err = tools.CopyDir(fs, configSrc, dotfileDest)
	} else {
		err = tools.CopyFile(fs, configSrc, dotfileDest)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Copied %s -> %s\n", configSrc, dotfileDest)
}
