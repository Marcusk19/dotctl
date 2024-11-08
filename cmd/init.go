package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Marcusk19/dotctl/tools"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(initCommand)
}

func copyExistingConfigs(programs []string, fs afero.Fs) {
	// takes list of programs and backs up configs for them
	destRoot := DotfilePath

	configRoot := ConfigPath
	for _, program := range programs {
		// TODO: do something here
		err := tools.CopyDir(fs, filepath.Join(configRoot, program), filepath.Join(destRoot, program))
		if err != nil {
			log.Fatalf("Problem copying %s", err.Error())
		}
	}
}

func createDotfileStructure(programs []string, fs afero.Fs) {
	// takes list of programs and creates dotfiles for them
	dotfileRoot := DotfilePath
	fmt.Printf("creating dotfile directory structure at %s\n", dotfileRoot)
	for _, program := range programs {
		if err := fs.MkdirAll(path.Join(dotfileRoot, program), os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Copy configs to dotfile directory",
	Long:  "Searches existing config directory for configs and then copies them to dotfile directory",
	Run:   runInitCommand,
}

func runInitCommand(cmd *cobra.Command, args []string) {
	fs := FileSystem
	// if user has passed a dotfile path flag need to add it to
	// viper's search path for a config file
	testing := viper.GetBool("testing")
	viper.AddConfigPath(filepath.Join(DotfilePath, "dotctl"))

	if viper.Get("testing") == true && fs.Name() != "MemMapFS" {
		log.Fatalf("wrong filesystem, got %s", fs.Name())
	}

	err := fs.MkdirAll(path.Join(DotfilePath, "dotctl"), 0755)
	if err != nil {
		log.Fatalf("Unable to create dotfile structure: %s", error.Error(err))
	}

	_, err = fs.Create(path.Join(DotfilePath, "dotctl/config.yml"))
	if err != nil {
		panic(fmt.Errorf("Unable to create config file %w", err))
	}

	if !testing {
		err = viper.WriteConfig()
		if err != nil && viper.Get("testing") != true {
			log.Fatalf("Unable to write config on init: %s\n", err)
		}

		_, err = git.PlainInit(DotfilePath, false)
		if err != nil {
			log.Fatal(err)
		}

		gitignoreContent := []byte(`
      # ignore dotctl config for individual installations
      dotctl/

      .DS_Store
      *.swp
      *.bak
      *.tmp
    `)

		err := afero.WriteFile(fs, filepath.Join(DotfilePath, ".gitignore"), gitignoreContent, 0644)

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully created dotfiles repository at %s\n", DotfilePath)
}
