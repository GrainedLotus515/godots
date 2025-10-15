# godots

A fast, interactive dotfiles installer built with Go and [Charmbracelet](https://github.com/charmbracelet).

## Features

- 🔗 **Symlink-based** - Edit dotfiles in place, changes sync to repo
- 💾 **Automatic backups** - Never lose your existing configs
- 🎨 **Interactive TUI** - Beautiful prompts with Charmbracelet Huh
- 🔄 **Auto-updates** - Integrates with pacman hooks
- 📦 **Multi-repo support** - Install from multiple dotfile repositories
- 🎯 **Selective install** - Choose which config groups to install
- 🪝 **Post-install hooks** - Run setup scripts after installation

## Installation

### From Source
```bash
git clone https://github.com/yourusername/godots
cd godots
go build -o godots
sudo mv godots /usr/local/bin/
```

### Quick Install Script
```bash
curl -sSL https://raw.githubusercontent.com/yourusername/godots/main/install.sh | bash
```

## Quick Start
```bash
# Install dotfiles from a repository
godots install https://github.com/yourusername/dots

# List installed repositories
godots list

# Update all dotfiles
godots update

# Update specific repository
godots update my-dots

# Uninstall a repository
godots uninstall my-dots

# Set up automatic updates on system upgrade
godots setup-hook
```

## Repository Structure

Your dotfiles repository should follow this structure:
```
dots/
├── config/           # Maps to ~/.config/
│   ├── nvim/
│   ├── zsh/
│   └── tmux/
├── local/            # Maps to ~/.local/
│   ├── bin/         # Scripts
│   └── share/       # Local data
├── home/             # Maps to ~/ (home root)
│   ├── .bashrc
│   ├── .zshrc
│   └── .gitconfig
└── hooks/            # Optional post-install scripts
    ├── nvim.sh
    └── zsh.sh
```

## How It Works

1. **Clone** - Repository is cloned to `~/.cache/godots/$reponame/`
2. **Scan** - Discovers dotfiles structure (config/, local/, home/)
3. **Select** - Interactive prompt to choose which groups to install
4. **Backup** - Existing files are backed up to `~/.godots.backup/TIMESTAMP/`
5. **Symlink** - Creates symlinks from cache to appropriate locations
6. **Hooks** - Optionally runs post-install scripts
7. **Manifest** - Tracks installation in `~/.config/godots/manifest.toml`

## Commands

### install

Install dotfiles from a git repository.
```bash
godots install <repository-url>
godots install https://github.com/user/dots
```

Options:
- `--auto` - Skip all prompts, install everything automatically

### list

List all installed dotfile repositories.
```bash
godots list
```

### update

Update dotfiles from git repository.
```bash
# Update all repositories
godots update

# Update specific repository
godots update my-dots
```

### uninstall

Remove a dotfiles repository.
```bash
godots uninstall <repo-name>
```

This will:
- Remove all symlinks
- Delete cached repository
- Update manifest
- Keep backups intact (manual cleanup)

### setup-hook

Install pacman hook for automatic updates.
```bash
godots setup-hook
```

After running this, your dotfiles will automatically update when you run `sudo pacman -Syu`.

### version

Print version information.
```bash
godots version
```

## Configuration

### Manifest File

Installation state is tracked in `~/.config/godots/manifest.toml`:
```toml
version = "1.0"

[repos.my-dots]
url = "https://github.com/user/dots"
cached_at = "/home/user/.cache/godots/my-dots"
installed_at = 2025-10-12T10:30:00Z
last_updated = 2025-10-12T10:30:00Z
installed_groups = ["nvim", "zsh", "tmux"]

  [repos.my-dots.symlinks]
  "/home/user/.config/nvim" = "config/nvim"
  "/home/user/.config/zsh" = "config/zsh"
  "/home/user/.local/bin/my-script" = "local/bin/my-script"
  "/home/user/.bashrc" = "home/.bashrc"
```

### Directory Structure
```
~/.cache/godots/       # Cloned repositories
~/.config/godots/      # Manifest and config
~/.godots.backup/      # Timestamped backups
```

## Hooks

Hooks are executable shell scripts in the `hooks/` directory of your dotfiles repository. They run after symlinks are created.

### Example Hook
```bash
#!/bin/bash
# hooks/nvim.sh

set -e

echo "Setting up Neovim..."

# Install lazy.nvim plugin manager
LAZY_PATH="${XDG_DATA_HOME:-$HOME/.local/share}/nvim/lazy/lazy.nvim"
if [ ! -d "$LAZY_PATH" ]; then
    git clone --filter=blob:none \
        https://github.com/folke/lazy.nvim.git \
        --branch=stable "$LAZY_PATH"
fi

# Install plugins
nvim --headless "+Lazy! sync" +qa

echo "✓ Neovim setup complete"
```

Make sure hooks are executable:
```bash
chmod +x hooks/*.sh
```

## Pacman Hook

After running `godots setup-hook`, a pacman hook is created at:
`~/.config/pacman/hooks/godots-update.hook`

This hook runs `godots update --auto` after every system upgrade.

## Examples

### Basic Installation
```bash
$ godots install https://github.com/user/dots

━━━ Installing Dotfiles ━━━
➜ Repository: https://github.com/user/dots
➜ Cloning repository...
✓ Cloned to /home/user/.cache/godots/dots
➜ Scanning dotfiles...
✓ Found 8 configuration groups

# Interactive selection prompt appears
# Select: nvim, zsh, tmux, git

➜ Checking for existing files...
⚠ Found 2 existing files
# Confirm backup? Yes

➜ Backing up existing files...
✓ Backed up to /home/user/.godots.backup/2025-10-12_10-30-00
➜ Creating symlinks...
✓ Created 4 symlinks
➜ Found 2 post-install hooks
# Run hooks? Yes
🔄 Running hook: nvim.sh
...
➜ Saving installation manifest...
✓ Manifest saved

━━━ Installation Complete! 🎉 ━━━
```

### Updating Dotfiles
```bash
$ godots update

➜ Updating my-dots...
Already up to date.
✓ my-dots updated
```

### Listing Installed Repos
```bash
$ godots list

━━━ Installed Dotfiles ━━━

📦 my-dots
   URL: https://github.com/user/dots
   Installed: 2025-10-12 10:30
   Groups: [nvim zsh tmux git]
```

## Development

### Building
```bash
go build -o godots
```

### Running Tests
```bash
go test ./...
```

### Project Structure
```
godots/
├── cmd/
│   └── main.go              # CLI commands
├── internal/
│   ├── config/              # Configuration management
│   ├── installer/           # Core installer logic
│   ├── manifest/            # TOML manifest handling
│   └── ui/                  # User interface (Huh + Lipgloss)
└── README.md
```

## Troubleshooting

### Symlink Creation Fails

**Issue**: Permission denied when creating symlinks

**Solution**: Ensure parent directories exist and you have write permissions
```bash
mkdir -p ~/.config ~/.local/bin
```

### Git Clone Fails

**Issue**: Repository not accessible

**Solution**: Check repository URL and SSH keys
```bash
# For private repos, ensure SSH key is added
ssh -T git@github.com
```

### Hooks Not Running

**Issue**: Hooks are discovered but don't execute

**Solution**: Make sure hooks are executable
```bash
chmod +x hooks/*.sh
```

### Backup Directory Full

**Issue**: `.godots.backup` taking up space

**Solution**: Old backups can be manually deleted
```bash
# List backups
ls -la ~/.godots.backup/

# Remove old backups (keep recent ones)
rm -rf ~/.godots.backup/2025-01-*
```

### Pacman Hook Not Working

**Issue**: Dotfiles don't update on system upgrade

**Solution**: Check hook is installed and permissions are correct
```bash
ls -la ~/.config/pacman/hooks/godots-update.hook
cat ~/.config/pacman/hooks/godots-update.hook
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [Charmbracelet Huh](https://github.com/charmbracelet/huh) for interactive prompts
- Uses [Charmbracelet Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- Manifest management with [BurntSushi TOML](https://github.com/BurntSushi/toml)
- CLI framework by [Cobra](https://github.com/spf13/cobra)

## Related Projects

- [chezmoi](https://github.com/twpayne/chezmoi) - More feature-complete dotfiles manager
- [GNU Stow](https://www.gnu.org/software/stow/) - Classic symlink farm manager
- [yadm](https://github.com/TheLocehiliosan/yadm) - Yet Another Dotfiles Manager
