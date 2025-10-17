package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/grainedlotus515/godotctl/internal/installer"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func PrintHeader(msg string) {
	fmt.Println(headerStyle.Render(fmt.Sprintf("\n━━━ %s ━━━", msg)))
}

func PrintSuccess(msg string) {
	fmt.Println(successStyle.Render("✓ " + msg))
}

func PrintError(msg string) {
	fmt.Println(errorStyle.Render("✗ " + msg))
}

func PrintInfo(msg string) {
	fmt.Println(infoStyle.Render("➜ " + msg))
}

func PrintWarning(msg string) {
	fmt.Println(warnStyle.Render("⚠ " + msg))
}

func PromptSelectGroups(groups []installer.DotfileGroup) ([]installer.DotfileGroup, error) {
	if len(groups) == 0 {
		return nil, fmt.Errorf("no groups available")
	}

	options := make([]huh.Option[string], len(groups))
	for i, group := range groups {
		options[i] = huh.NewOption(
			fmt.Sprintf("%s (%s)", group.Name, group.Target),
			group.Name,
		)
	}

	var selected []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select configuration groups to install").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	// Filter groups based on selection
	var result []installer.DotfileGroup
	for _, group := range groups {
		for _, sel := range selected {
			if group.Name == sel {
				result = append(result, group)
				break
			}
		}
	}

	return result, nil
}

func PromptConfirm(message string) (bool, error) {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirm, nil
}
