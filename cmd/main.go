package main

import (
	"fmt"
	"os"

	"github.com/GrainedLotus515/dotfiles-installer/internal/config"
	"github.com/GrainedLotus515/dotfiles-installer/internal/installer"
	"github.com/GrainedLotus515/dotfiles-installer/internal/manifest"
	"github.com/GrainedLotus515/dotfiles-installer/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	auto    bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dotfiles-installer",
	Short: "A dotfiles installer for CachyOS",
	Long:  `Manage your dotfiles with symlinks, backups, and hooks.`,
}

var installCmd = &cobra.Command{
	Use:   "install [repository-url]",
	Short: "Install dotfiles from a git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]

		ui.PrintHeader("Installing Dotfiles")
		ui.PrintInfo(fmt.Sprintf("Repository: %s", repoURL))

		// Initialize configuration
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// Initialize installer
		inst := installer.New(cfg)

		// Clone repository
		ui.PrintInfo("Cloning repository...")
		repoPath, repoName, err := inst.Clone(repoURL)
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		ui.PrintSuccess(fmt.Sprintf("Cloned to %s", repoPath))

		// Scan dotfiles structure
		ui.PrintInfo("Scanning dotfiles...")
		groups, err := inst.Scan(repoPath)
		if err != nil {
			return fmt.Errorf("failed to scan dotfiles: %w", err)
		}
		ui.PrintSuccess(fmt.Sprintf("Found %d configuration groups", len(groups)))

		// Show selection prompt (unless auto mode)
		var selectedGroups []installer.DotfileGroup
		if auto {
			selectedGroups = groups
		} else {
			selectedGroups, err = ui.PromptSelectGroups(groups)
			if err != nil {
				return fmt.Errorf("selection cancelled: %w", err)
			}
		}

		if len(selectedGroups) == 0 {
			ui.PrintInfo("No groups selected, exiting")
			return nil
		}

		// Check for conflicts and backup
		ui.PrintInfo("Checking for existing files...")
		conflicts, err := inst.CheckConflicts(selectedGroups)
		if err != nil {
			return fmt.Errorf("failed to check conflicts: %w", err)
		}

		if len(conflicts) > 0 {
			ui.PrintWarning(fmt.Sprintf("Found %d existing files", len(conflicts)))

			if !auto {
				confirm, err := ui.PromptConfirm("Backup existing files and continue?")
				if err != nil || !confirm {
					return fmt.Errorf("installation cancelled")
				}
			}

			ui.PrintInfo("Backing up existing files...")
			backupDir, err := inst.Backup(conflicts)
			if err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
			ui.PrintSuccess(fmt.Sprintf("Backed up to %s", backupDir))
		}

		// Create symlinks
		ui.PrintInfo("Creating symlinks...")
		symlinks, err := inst.CreateSymlinks(selectedGroups)
		if err != nil {
			return fmt.Errorf("failed to create symlinks: %w", err)
		}
		ui.PrintSuccess(fmt.Sprintf("Created %d symlinks", len(symlinks)))

		// Discover hooks
		hooks, err := inst.DiscoverHooks(repoPath)
		if err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to discover hooks: %v", err))
		} else if len(hooks) > 0 {
			ui.PrintInfo(fmt.Sprintf("Found %d post-install hooks", len(hooks)))

			if !auto {
				runHooks, err := ui.PromptConfirm("Run post-install hooks?")
				if err == nil && runHooks {
					if err := inst.RunHooks(hooks, false); err != nil {
						ui.PrintWarning(fmt.Sprintf("Some hooks failed: %v", err))
					}
				}
			} else {
				if err := inst.RunHooks(hooks, true); err != nil {
					ui.PrintWarning(fmt.Sprintf("Some hooks failed: %v", err))
				}
			}
		}

		// Save manifest
		ui.PrintInfo("Saving installation manifest...")
		man := manifest.New(cfg.ManifestPath)
		if err := man.AddRepo(repoName, repoURL, repoPath, selectedGroups, symlinks); err != nil {
			return fmt.Errorf("failed to save manifest: %w", err)
		}
		ui.PrintSuccess("Manifest saved")

		ui.PrintHeader("Installation Complete! 🎉")
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed dotfile repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		man := manifest.New(cfg.ManifestPath)
		repos, err := man.Load()
		if err != nil {
			return fmt.Errorf("failed to load manifest: %w", err)
		}

		if len(repos) == 0 {
			ui.PrintInfo("No dotfiles installed")
			return nil
		}

		ui.PrintHeader("Installed Dotfiles")
		for name, repo := range repos {
			fmt.Printf("\n📦 %s\n", name)
			fmt.Printf("   URL: %s\n", repo.URL)
			fmt.Printf("   Installed: %v\n", repo.InstalledAt.Format("2006-01-02 15:04"))
			fmt.Printf("   Groups: %v\n", repo.InstalledGroups)
		}

		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [repo-name]",
	Short: "Update installed dotfiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		man := manifest.New(cfg.ManifestPath)
		repos, err := man.Load()
		if err != nil {
			return fmt.Errorf("failed to load manifest: %w", err)
		}

		if len(repos) == 0 {
			ui.PrintInfo("No dotfiles installed")
			return nil
		}

		// If specific repo provided, update only that one
		if len(args) > 0 {
			repoName := args[0]
			if repo, exists := repos[repoName]; exists {
				return updateRepo(cfg, repoName, repo)
			}
			return fmt.Errorf("repository '%s' not found", repoName)
		}

		// Update all repos
		for name, repo := range repos {
			if err := updateRepo(cfg, name, repo); err != nil {
				ui.PrintWarning(fmt.Sprintf("Failed to update %s: %v", name, err))
			}
		}

		return nil
	},
}

func updateRepo(cfg *config.Config, name string, repo manifest.RepoConfig) error {
	ui.PrintInfo(fmt.Sprintf("Updating %s...", name))

	inst := installer.New(cfg)

	// Git pull in cached repo
	if err := inst.Update(repo.CachedAt); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("%s updated", name))
	return nil
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [repo-name]",
	Short: "Uninstall dotfiles repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoName := args[0]

		cfg, err := config.New()
		if err != nil {
			return err
		}

		man := manifest.New(cfg.ManifestPath)
		repos, err := man.Load()
		if err != nil {
			return err
		}

		repo, exists := repos[repoName]
		if !exists {
			return fmt.Errorf("repository '%s' not found", repoName)
		}

		ui.PrintWarning(fmt.Sprintf("This will remove %d symlinks from %s", len(repo.Symlinks), repoName))

		confirm, err := ui.PromptConfirm("Continue with uninstall?")
		if err != nil || !confirm {
			return fmt.Errorf("uninstall cancelled")
		}

		inst := installer.New(cfg)

		// Remove symlinks
		ui.PrintInfo("Removing symlinks...")
		if err := inst.RemoveSymlinks(repo.Symlinks); err != nil {
			return err
		}

		// Remove cached repo
		ui.PrintInfo("Removing cached repository...")
		if err := os.RemoveAll(repo.CachedAt); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to remove cache: %v", err))
		}

		// Update manifest
		if err := man.RemoveRepo(repoName); err != nil {
			return err
		}

		ui.PrintSuccess("Uninstalled successfully")
		return nil
	},
}

var setupHookCmd = &cobra.Command{
	Use:   "setup-hook",
	Short: "Install pacman hook for auto-updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		inst := installer.New(cfg)

		ui.PrintInfo("Installing pacman hook...")
		if err := inst.SetupPacmanHook(); err != nil {
			return err
		}

		ui.PrintSuccess("Pacman hook installed at ~/.config/pacman/hooks/dotfiles-update.hook")
		ui.PrintInfo("Dotfiles will now auto-update on system upgrades")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotfiles-installer v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(setupHookCmd)
	rootCmd.AddCommand(versionCmd)

	installCmd.Flags().BoolVar(&auto, "auto", false, "Automatic mode (no prompts)")
}
