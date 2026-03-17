package styles

import "github.com/charmbracelet/lipgloss"

const (
	InfoColor     = "#5dd2fc"
	ErrorColor    = "#ba0b0b"
	SelectedColor = "#FC895F"
	FocusedColor  = "#FCBD5F"
	DefaultColor  = "#5dd2fc"
	BlurredColor  = "240"
)

var (
	InfoStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(InfoColor)).Bold(true)
	ErrorStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(ErrorColor)).Bold(true)
	SelectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(SelectedColor)).Bold(true)
	DefaultStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(DefaultColor)).Bold(true)
	FocusedStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(FocusedColor)).Bold(true)
	BlurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(BlurredColor))
	FocusedButton = SelectedStyle.Render("[ Submit ]")
	BlurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color(BlurredColor)).Render("[ Submit ]")
)
