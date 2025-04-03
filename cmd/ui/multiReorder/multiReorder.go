package multiReorder

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logoMerge = `
 __  __                  
|  \/  |___ _ _ __ _ ___ 
| |\/| / -_) '_/ _` + "`" + ` / -_)
|_|  |_\___|_| \__, \___|
               |___/     
`

var (
	defaultStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#5dd2fc")).Bold(true)
	focusedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FCBD5F")).Bold(true)
)

type Tmodel struct {
	pdfs   []string
	cursor int
	logo   string
	Quit   bool
}

func MultiReorderModel(pdfs []string, logo string) Tmodel {
	return Tmodel{
		pdfs: pdfs,
		logo: logo,
	}
}

func (m Tmodel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Tmodel) GetOrderedPdfs() []string {
	return m.pdfs
}

func (m Tmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quit = true
			return m, tea.Quit

		case "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "j":
			if m.cursor < len(m.pdfs)-1 {
				m.cursor++
			}

		case "up":
			if m.cursor > 0 {
				m.pdfs[m.cursor], m.pdfs[m.cursor-1] = m.pdfs[m.cursor-1], m.pdfs[m.cursor]
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.pdfs)-1 {
				m.pdfs[m.cursor], m.pdfs[m.cursor+1] = m.pdfs[m.cursor+1], m.pdfs[m.cursor]
				m.cursor++
			}

		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Tmodel) View() string {
	var b strings.Builder
	b.WriteString((defaultStyle.Render(logoMerge)))

	fmt.Fprint(&b, "\n\n")
	b.WriteString(focusedStyle.Render("Reorder the PDFs:"))
	fmt.Fprint(&b, "\n")
	b.WriteString(focusedStyle.Render("Navigate using the 'j/k' keys, reorder using 'up/down' keys."))
	fmt.Fprint(&b, "\n\n")

	for i, choice := range m.pdfs {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			choice = focusedStyle.Render(choice)
		}

		choice = defaultStyle.Render(choice)

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	b.WriteString(focusedStyle.Render("\nPress enter to confirm, esc to quit."))
	fmt.Fprint(&b, "\n\n")

	return b.String()
}

func MultiReorderInteractive(pdfs []string, logo string) (reorderedPdfs []string, quit bool, err error) {
	r := tea.NewProgram(MultiReorderModel(pdfs, logo))
	result, err := r.Run()
	if err != nil {
		return nil, false, err
	}

	model := result.(Tmodel)
	if model.Quit {
		return nil, true, fmt.Errorf("operation canceled")
	}

	selectedPdfs := model.GetOrderedPdfs()

	return selectedPdfs, false, nil
}
