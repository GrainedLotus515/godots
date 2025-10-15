package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (i *Installer) CheckConflicts(groups []DotfileGroup) ([]string, error) {
	var conflicts []string

	for _, group := range groups {
		if _, err := os.Lstat(group.Target); err == nil {
			conflicts = append(conflicts, group.Target)
		}
	}

	return conflicts, nil
}

func (i *Installer) Backup(paths []string) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupDir := filepath.Join(i.cfg.BackupDir, timestamp)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	for _, path := range paths {
		// Get relative path from home
		relPath, err := filepath.Rel(i.cfg.HomeDir, path)
		if err != nil {
			return "", err
		}

		backupPath := filepath.Join(backupDir, relPath)
		backupParent := filepath.Dir(backupPath)

		// Create parent directories
		if err := os.MkdirAll(backupParent, 0755); err != nil {
			return "", err
		}

		// Move file to backup
		if err := os.Rename(path, backupPath); err != nil {
			// If rename fails, try copy
			if err := copyFile(path, backupPath); err != nil {
				return "", fmt.Errorf("failed to backup %s: %w", path, err)
			}
			os.Remove(path)
		}
	}

	return backupDir, nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
