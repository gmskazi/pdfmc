package styles

import "github.com/charmbracelet/lipgloss"

const (
	InfoColor     = "#5dd2fc"
	ErrorColor    = "#ba0b0b"
	SelectedColor = "#FC895F"
)

var (
	InfoStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(InfoColor)).Bold(true)
	ErrorStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(ErrorColor)).Bold(true)
	SelectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(SelectedColor)).Bold(true)
)
