# Dotctl

A cli tool to manage your dotfiles
## About
Dotctl is a tool to help you easily manage your dotfiles and sync them across separate machines using
git. It aims to abstract away the manual effort of symlinking your dotfiles to config directories and
updating them with git.

## Installation

### Build From Source
_Prerequisites_
- [go](https://go.dev/doc/install)

clone the repo and run script to build binary and copy it to your path

```sh
git clone https://github.com/Marcusk19/dotctl.git
cd dotcl
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
# sync
dotctl sync -r <remote-repo>
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

### Docker
You can also run the docker container that will set up a barebones shell for you
to test dotctl with 

```sh
docker built -d dotctl-dev .
docker run -it dotctl-dev
dotctl status
```

