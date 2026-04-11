# dotctl — Fresh Machine Bootstrap Roadmap

## Goal

A single command on a fresh machine should clone a dotfiles repo and fully apply it:

```bash
dotctl apply https://github.com/user/dotfiles.git
```

No manual git clone. No rebuilding from source. No multi-step dance.

---

## Current State Summary

| Capability | Status |
|---|---|
| `init` / `add` / `link` core loop | ✅ Working |
| Symlink creation | ✅ Working |
| Config stored at `~/dotfiles/dotctl/config.yaml` | ✅ Working |
| `sync` command (pull + push) | ⚠️ Unstable — README warns against it |
| Bootstrap from existing repo | ❌ Missing |
| Conflict detection on `link` | ❌ Missing |
| Idempotent `link` (safe to re-run) | ❌ Missing |
| One-liner install script | ❌ Missing |
| Templating / per-machine config | ❌ Not in scope (yet) |
| Secret management | ❌ Not in scope (yet) |

---

## Phase 1 — Fix the Foundation (prerequisite for everything)

These are bugs/gaps that block safe automation.

### 1.1 Make `link` idempotent

**Problem:** Running `link` a second time likely fails or clobbers without warning if a target path already exists.

**Fix:** Before creating each symlink, check the target:
- If missing → create symlink (happy path)
- If already a symlink pointing to the correct source → skip, log "already linked"
- If already a symlink pointing somewhere else → warn, skip unless `--overwrite` flag passed
- If a real file/directory exists → back it up to `<path>.dotctl.bak`, then link (or skip with `--no-backup`)

**Files to touch:** `cmd/link.go`

**Acceptance criteria:** Running `dotctl link` twice on a clean setup produces no errors and no duplicate symlinks.

---

### 1.2 Fix or remove `sync`

**Problem:** The README actively warns against `sync`. It should either be fixed or removed so it doesn't cause data loss for new users.

**Decision to make (pick one):**
- **Option A — Fix it:** Proper git pull → detect conflicts → commit changed files → push. Use `go-git` for in-process git operations instead of shelling out.
- **Option B — Remove it:** Drop `sync`, document that users should manage the `~/dotfiles` directory as a normal git repo. Add a note in README pointing to how to do `git pull` and then `dotctl link`.

Recommendation: **Option B** now, **Option A** later in Phase 3. It's safer to remove a broken command than ship a half-working one.

---

### 1.3 Validate that config round-trips cleanly

**Problem:** The tracked file list lives in `~/dotfiles/dotctl/config.yaml`. On a fresh machine, this file needs to already exist in the cloned repo for bootstrap to work. Need to confirm this file is actually committed and not gitignored.

**Fix:** Add a note in `.gitignore` explicitly *not* ignoring `dotctl/config.yaml`. Add a check in `init` that warns if the config file would be gitignored.

---

## Phase 2 — The Bootstrap Command (core deliverable)

### 2.1 `dotctl apply <repo-url>`

This is the main feature. It should:

1. Check if `~/dotfiles` already exists
   - If yes and it's a git repo → `git pull` (or prompt user)
   - If yes but not a git repo → error with clear message
   - If no → `git clone <repo-url> ~/dotfiles`
2. Read `~/dotfiles/dotctl/config.yaml` (fail clearly if not found — tells user their repo isn't set up for dotctl)
3. Run the equivalent of `dotctl link` with idempotent behavior from Phase 1.1
4. Print a summary: N linked, M skipped, K backed up

**Flags:**
- `--overwrite` — overwrite existing files instead of backing up
- `--dry-run` — print what would happen without touching the filesystem
- `--no-backup` — skip backups, just skip conflicts

**Files to touch:** New `cmd/apply.go`

**Acceptance criteria:** On a machine with nothing but Go installed and git in PATH, running `dotctl apply https://github.com/user/dotfiles.git` produces a fully linked dotfiles setup.

---

### 2.2 One-liner install + apply script

The bootstrap UX needs to work before `dotctl` itself is installed. Options:

**Option A — `go install` (simplest):**
```bash
go install github.com/Marcusk19/dotctl@latest && dotctl apply https://github.com/user/dotfiles.git
```
Requires Go on the machine. Works great for developer setups.

**Option B — Shell script + prebuilt binary:**
```bash
curl -fsSL https://raw.githubusercontent.com/Marcusk19/dotctl/main/install.sh | bash -s -- https://github.com/user/dotfiles.git
```
The script: detects OS/arch → downloads the correct binary from GitHub releases → runs `dotctl apply <repo>`. No Go required.

Recommendation: **ship both**. `go install` for devs, the curl script for fresh machines where Go may not be present. goreleaser is already set up, so binary releases exist — just need the install script.

**Files to add:** `install.sh` in repo root

---

## Phase 3 — Quality of Life (after core bootstrap works)

### 3.1 Fix `sync` properly

With `go-git` as a dependency (or shelling to git with proper error handling):
- `dotctl sync` → pull remote changes, re-run link, push any local changes
- Detect and surface merge conflicts clearly instead of silently failing

---

### 3.2 `dotctl status`

Shows current state of all tracked files:
```
nvim        → ~/.config/nvim         [linked ✓]
zshrc       → ~/.zshrc               [linked ✓]
kitty       → ~/.config/kitty        [CONFLICT — real file exists]
tmux.conf   → ~/.tmux.conf           [broken symlink]
```

---

### 3.3 `dotctl remove <path>`

Removes a tracked file from the manifest and optionally replaces the symlink with the real file.

---

### 3.4 `dotctl list`

Prints all tracked entries from config — useful for auditing before running `apply` on a new machine.

---

## Phase 4 — Advanced (optional, post-stabilization)

| Feature | Notes |
|---|---|
| Per-machine profiles | Tag entries in config with `profiles: [work, home]`; pass `--profile` to `apply` |
| Templating | Go `text/template` over files with a `.tmpl` extension — render before linking |
| Secret redaction | Warn when a tracked file contains patterns that look like secrets (API keys, tokens) |
| XDG-aware paths | Auto-resolve `$XDG_CONFIG_HOME` instead of hardcoding `~/.config` |

---

## Implementation Order

```
Phase 1.1 → idempotent link        (1–2 days)
Phase 1.2 → remove/fix sync        (0.5 days)
Phase 1.3 → config gitignore audit (0.5 days)
Phase 2.1 → dotctl apply command   (2–3 days)
Phase 2.2 → install.sh script      (1 day)
Phase 3.x → status, list, remove   (2–3 days)
Phase 4.x → profiles, templating   (future)
```

---

## Definition of Done (for "point at a repo and go")

- [ ] `install.sh` downloads the correct binary for current OS/arch
- [ ] `dotctl apply <url>` clones repo if not present, reads config, links all tracked files
- [ ] Running `apply` twice produces no errors
- [ ] Existing files are backed up, not silently overwritten
- [ ] `--dry-run` works and shows exactly what would happen
- [ ] README updated with the new one-liner bootstrap flow
- [ ] At least one end-to-end test that: clones a test dotfiles repo → runs apply → asserts symlinks exist
