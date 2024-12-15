# Dotctl

dotfile management
## About
Dotctl is a tool to help you easily manage your dotfiles and sync them across separate machines using
git. It creates a `dotfiles` subdirectory in the user's `$HOME` and provides simple commands to add
and symlink config files/directories to the central `dotfiles` directory.


## Installation

### Build From Source
_Prerequisites_
- [go](https://go.dev/doc/install)

clone the repo and run script to build binary and copy it to your path

```sh
git clone https://github.com/Marcusk19/dotctl.git
cd dotctl
make install
```

## Usage

```bash
# init sets up the config file and directory to hold all dotfiles
dotctl init
# add a config directory for dotctl to track
dotctl add ~/.config/nvim
# create symlinks
dotctl link
```
### Syncing to git
_Warning: using the sync command can have some unexpected behavior, currently the recommendation
is to manually track the dotfiles with git_

dotctl comes with a `sync` command that performs the following operations for the dotfiles directory:

1. pulls changes from configured upstream git repo 
2. commits and pushes any changes detected in the dotfile repo

set the upstream repo using the `-r` flag or manually edit the config at `$HOME/dotfiles/dotctl/config.yaml`

example usage:
```
dotctl sync -r https://github.com/example/dotfiles.git
```
