package textInputs

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/gmskazi/pdfmergecrypt/cmd/ui/textinputs"
)

var (
	defaultStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#5dd2fc")).Bold(true)
	helpStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#F1F0E9")).Bold(true)
	focusedStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FCBD5F")).Bold(true)
	errorStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#ba0b0b")).Bold(true)
	noStyle       = lipgloss.NewStyle()
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Tmodel struct {
	focusIndex int
	inputs     []textinputs.Tmodel
	Quit       bool
}
