package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (i *Installer) Clone(repoURL string) (string, string, error) {
	// Extract repo name from URL
	repoName := extractRepoName(repoURL)
	repoPath := filepath.Join(i.cfg.CacheDir, repoName)

	// Check if already cloned
	if _, err := os.Stat(repoPath); err == nil {
		// Already exists, just return path
		return repoPath, repoName, nil
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", repoURL, repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("git clone failed: %w", err)
	}

	return repoPath, repoName, nil
}

func (i *Installer) Update(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return nil
}

func extractRepoName(url string) string {
	// Extract repo name from git URL
	// https://github.com/user/repo.git -> repo
	// https://github.com/user/repo -> repo

	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return "dotfiles"
	}

	name := parts[len(parts)-1]
	name = strings.TrimSuffix(name, ".git")

	if name == "" {
		return "dotfiles"
	}

	return name
}
