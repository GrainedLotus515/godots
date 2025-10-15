package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	HomeDir      string
	CacheDir     string
	ConfigDir    string
	ManifestPath string
	BackupDir    string
}

func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "godots")
	configDir := filepath.Join(homeDir, ".config", "godots")
	manifestPath := filepath.Join(configDir, "manifest.toml")
	backupDir := filepath.Join(homeDir, ".godots.backup")

	// Ensure directories exist
	for _, dir := range []string{cacheDir, configDir, backupDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	return &Config{
		HomeDir:      homeDir,
		CacheDir:     cacheDir,
		ConfigDir:    configDir,
		ManifestPath: manifestPath,
		BackupDir:    backupDir,
	}, nil
}
