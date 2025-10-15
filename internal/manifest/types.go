package manifest

import "time"

type Manifest struct {
	Version string                `toml:"version"`
	Repos   map[string]RepoConfig `toml:"repos"`
}

type RepoConfig struct {
	URL             string            `toml:"url"`
	CachedAt        string            `toml:"cached_at"`
	InstalledAt     time.Time         `toml:"installed_at"`
	LastUpdated     time.Time         `toml:"last_updated"`
	InstalledGroups []string          `toml:"installed_groups"`
	Symlinks        map[string]string `toml:"symlinks"`
}
