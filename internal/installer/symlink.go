package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

func (i *Installer) CreateSymlinks(groups []DotfileGroup) (map[string]string, error) {
	symlinks := make(map[string]string)

	for _, group := range groups {
		// Ensure parent directory exists
		parent := filepath.Dir(group.Target)
		if err := os.MkdirAll(parent, 0755); err != nil {
			return nil, fmt.Errorf("failed to create parent dir %s: %w", parent, err)
		}

		// Create symlink
		if err := os.Symlink(group.Source, group.Target); err != nil {
			return nil, fmt.Errorf("failed to create symlink %s -> %s: %w", group.Target, group.Source, err)
		}

		symlinks[group.Target] = group.Source
	}

	return symlinks, nil
}

func (i *Installer) RemoveSymlinks(symlinks map[string]string) error {
	for target := range symlinks {
		// Verify it's actually a symlink before removing
		info, err := os.Lstat(target)
		if err != nil {
			continue // Already gone
		}

		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove symlink %s: %w", target, err)
			}
		}
	}

	return nil
}
