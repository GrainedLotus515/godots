package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Hook struct {
	Name string
	Path string
}

func (i *Installer) DiscoverHooks(repoPath string) ([]Hook, error) {
	hooksDir := filepath.Join(repoPath, "hooks")

	// Check if hooks directory exists
	if _, err := os.Stat(hooksDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(hooksDir)
	if err != nil {
		return nil, err
	}

	var hooks []Hook
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if executable
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.Mode()&0111 != 0 { // Has execute permission
			hooks = append(hooks, Hook{
				Name: entry.Name(),
				Path: filepath.Join(hooksDir, entry.Name()),
			})
		}
	}

	return hooks, nil
}

func (i *Installer) RunHooks(hooks []Hook, silent bool) error {
	for _, hook := range hooks {
		if !silent {
			fmt.Printf("🔄 Running hook: %s\n", hook.Name)
		}

		cmd := exec.Command("bash", hook.Path)
		cmd.Dir = i.cfg.HomeDir

		if !silent {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %s failed: %w", hook.Name, err)
		}
	}

	return nil
}
