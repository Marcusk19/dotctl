# Dotctl

A cli tool to manage your dotfiles
## About
Dotctl is a tool to help you easily manage your dotfiles and sync them across separate machines using
git. It aims to abstract away the manual effort of symlinking your dotfiles to config directories and
updating them with git.

## Installation
- TBD

## Usage

_Note: at the moment the sync feature has some bugs that need to be fixed, to_
_avoid breaking your configs it is recommended to manually use git for the time being_

```bash
# init sets up the config file and directory to hold all dotfiles
dotctl init
# add a config directory for dotctl to track
dotctl add ~/.config/nvim
# create symlinks
dotctl link
```

## Development
It's preferable to create a temporary directory and copy your system's config
directory over to avoid making undesirable changes to your system.
A couple of useful makefile scripts exist to set up and tear down this.
It will create a testing directory in `./tmp/config` and copy your system configs
over.

```bash
make sandbox # creates the directory and copies over from ~/.config
make clean # removes directory
```


