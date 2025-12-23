package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SourceType string

const (
	SourceTypeRemote   SourceType = "remote"    // Remote git repository
	SourceTypeLocalGit SourceType = "local-git" // Local git repository
	SourceTypeLocalDir SourceType = "local-dir" // Local directory (non-git)
)

// Clone handles cloning/copying a repository from various sources
func (i *Installer) Clone(source string) (string, string, SourceType, error) {
	sourceType := detectSourceType(source)
	repoName := extractRepoName(source)
	repoPath := filepath.Join(i.cfg.CacheDir, repoName)

	// Check if already exists in cache
	if _, err := os.Stat(repoPath); err == nil {
		// Already exists, just return path
		return repoPath, repoName, sourceType, nil
	}

	switch sourceType {
	case SourceTypeRemote:
		// Clone remote git repository
		cmd := exec.Command("git", "clone", source, repoPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", "", "", fmt.Errorf("git clone failed: %w", err)
		}

	case SourceTypeLocalGit, SourceTypeLocalDir:
		// Copy local directory to cache
		if err := copyDir(source, repoPath); err != nil {
			return "", "", "", fmt.Errorf("failed to copy local directory: %w", err)
		}
	}

	return repoPath, repoName, sourceType, nil
}

// Update handles updating a repository based on its source type
func (i *Installer) Update(repoPath string, sourceType SourceType) error {
	switch sourceType {
	case SourceTypeRemote, SourceTypeLocalGit:
		// Check if it's a git repository
		gitDir := filepath.Join(repoPath, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// Git pull for remote repos or local git repos with remotes
			cmd := exec.Command("git", "-C", repoPath, "pull")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				// If pull fails (no remote configured for local git), that's okay
				if sourceType == SourceTypeLocalGit {
					return nil
				}
				return fmt.Errorf("git pull failed: %w", err)
			}
		}

	case SourceTypeLocalDir:
		// For local directories, we don't update (they're snapshots)
		return nil
	}

	return nil
}

// detectSourceType determines if the source is a remote URL, local git repo, or local directory
func detectSourceType(source string) SourceType {
	// Check if it's a URL (http://, https://, git@, etc.)
	if strings.Contains(source, "://") || strings.HasPrefix(source, "git@") {
		return SourceTypeRemote
	}

	// Check if it's a local path
	absPath, err := filepath.Abs(source)
	if err != nil {
		// If we can't resolve the path, assume it's a remote URL
		return SourceTypeRemote
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		// Path doesn't exist, assume it's a remote URL
		return SourceTypeRemote
	}

	// Check if it's a git repository
	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return SourceTypeLocalGit
	}

	// It's a local directory but not a git repo
	return SourceTypeLocalDir
}

func extractRepoName(source string) string {
	// Extract repo name from source
	// https://github.com/user/repo.git -> repo
	// https://github.com/user/repo -> repo
	// /path/to/dotfiles -> dotfiles

	// For local paths, use the directory name
	if !strings.Contains(source, "://") && !strings.HasPrefix(source, "git@") {
		absPath, err := filepath.Abs(source)
		if err == nil {
			source = absPath
		}
	}

	parts := strings.Split(source, "/")
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

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file with mode preservation
			if err := copyFilePreserveMode(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFilePreserveMode copies a single file while preserving its permissions
func copyFilePreserveMode(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
