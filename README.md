# Dotctl

dotfile management

## About

Dotctl is a tool to help you easily manage your dotfiles and sync them across separate machines using
git. It creates a `dotfiles` subdirectory in the user's `$HOME` and provides simple commands to add
and symlink config files/directories to the central `dotfiles` directory.

## Quick Start (Bootstrap a fresh machine)

**Option 1 — Shell script (no Go required):**

```bash
curl -fsSL https://raw.githubusercontent.com/Marcusk19/dotctl/main/install.sh | bash -s -- https://github.com/your-user/dotfiles.git
```

**Option 2 — go install:**

```bash
go install github.com/Marcusk19/dotctl@latest
dotctl apply https://github.com/your-user/dotfiles.git
```

Both methods clone your dotfiles repo to `~/dotfiles` and link all tracked configs automatically.

## Installation

### go install

```sh
go install github.com/Marcusk19/dotctl@latest
```

### Build From Source

_Prerequisites_
- [go](https://go.dev/doc/install)

Clone the repo and run the script to build the binary and copy it to your path:

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

# create symlinks (idempotent, safe to re-run)
dotctl link

# bootstrap dotfiles on a fresh machine
dotctl apply https://github.com/your-user/dotfiles.git
```

### Commands

| Command | Description |
|---------|-------------|
| `dotctl init` | Set up a new dotfiles repo and config |
| `dotctl add <path>` | Track a config file or directory |
| `dotctl link` | Create symlinks for all tracked configs |
| `dotctl apply <repo-url>` | Clone a dotfiles repo and link everything (bootstrap) |

### Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Show what would be done without making changes |
| `--overwrite` | Overwrite existing files when linking |
| `--no-backup` | Skip creating backups of existing files |
