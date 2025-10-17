package manifest

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/grainedlotus515/godotctl/internal/installer"
)

type Manager struct {
	path string
}

func New(path string) *Manager {
	return &Manager{path: path}
}

func (m *Manager) Load() (map[string]RepoConfig, error) {
	var manifest Manifest

	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		return make(map[string]RepoConfig), nil
	}

	if _, err := toml.DecodeFile(m.path, &manifest); err != nil {
		return nil, err
	}

	if manifest.Repos == nil {
		manifest.Repos = make(map[string]RepoConfig)
	}

	return manifest.Repos, nil
}

func (m *Manager) Save(repos map[string]RepoConfig) error {
	manifest := Manifest{
		Version: "1.0",
		Repos:   repos,
	}

	f, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(manifest)
}

func (m *Manager) AddRepo(name, url, cachedAt string, groups []installer.DotfileGroup, symlinks map[string]string) error {
	repos, err := m.Load()
	if err != nil {
		return err
	}

	groupNames := make([]string, len(groups))
	for i, g := range groups {
		groupNames[i] = g.Name
	}

	repos[name] = RepoConfig{
		URL:             url,
		CachedAt:        cachedAt,
		InstalledAt:     time.Now(),
		LastUpdated:     time.Now(),
		InstalledGroups: groupNames,
		Symlinks:        symlinks,
	}

	return m.Save(repos)
}

func (m *Manager) RemoveRepo(name string) error {
	repos, err := m.Load()
	if err != nil {
		return err
	}

	delete(repos, name)
	return m.Save(repos)
}

func (m *Manager) UpdateRepo(name string) error {
	repos, err := m.Load()
	if err != nil {
		return err
	}

	if repo, exists := repos[name]; exists {
		repo.LastUpdated = time.Now()
		repos[name] = repo
		return m.Save(repos)
	}

	return nil
}
