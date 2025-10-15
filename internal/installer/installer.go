package installer

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/GrainedLotus515/godots/internal/config"
)

type Installer struct {
    cfg *config.Config
}

func New(cfg *config.Config) *Installer {
    return &Installer{cfg: cfg}
}

func (i *Installer) SetupPacmanHook() error {
    hookDir := filepath.Join(i.cfg.HomeDir, ".config", "pacman", "hooks")
    if err := os.MkdirAll(hookDir, 0755); err != nil {
        return err
    }
    
    hookPath := filepath.Join(hookDir, "dotfiles-update.hook")
    
    hookContent := `[Trigger]
Operation = Upgrade
Type = Package
Target = *

[Action]
Description = Updating dotfiles repositories...
When = PostTransaction
Exec = /usr/local/bin/dotfiles-installer update --auto
`
    
    if err := os.WriteFile(hookPath, []byte(hookContent), 0644); err != nil {
        return fmt.Errorf("failed to write hook file: %w", err)
    }
    
    return nil
}