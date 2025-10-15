package installer

import (
	"os"
	"path/filepath"
)

type DotfileGroup struct {
	Name   string
	Source string
	Target string
	Files  []string
}

type PathMapping struct {
	SourceDir string
	TargetDir string
}

func (i *Installer) Scan(repoPath string) ([]DotfileGroup, error) {
	var groups []DotfileGroup

	mappings := []PathMapping{
		{SourceDir: "config", TargetDir: filepath.Join(i.cfg.HomeDir, ".config")},
		{SourceDir: "local", TargetDir: filepath.Join(i.cfg.HomeDir, ".local")},
		{SourceDir: "home", TargetDir: i.cfg.HomeDir},
	}

	for _, mapping := range mappings {
		sourcePath := filepath.Join(repoPath, mapping.SourceDir)

		// Check if directory exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			continue
		}

		// Scan for subdirectories/files
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			// Skip hooks directory
			if entry.Name() == "hooks" {
				continue
			}

			source := filepath.Join(sourcePath, entry.Name())
			target := filepath.Join(mapping.TargetDir, entry.Name())

			group := DotfileGroup{
				Name:   entry.Name(),
				Source: source,
				Target: target,
				Files:  []string{entry.Name()},
			}

			groups = append(groups, group)
		}
	}

	return groups, nil
}
