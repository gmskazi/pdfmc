package multiSelect

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

const logoEncrypt = `
 ___                       _   
| __|_ _  __ _ _ _  _ _ __| |_ 
| _|| ' \/ _| '_| || | '_ \  _|
|___|_||_\__|_|  \_, | .__/\__|
                 |__/|_|       
`

var (
	defaultStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#5dd2fc")).Bold(true)
	focusedStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FCBD5F")).Bold(true)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FC895F")).Bold(true)
	errorStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#ba0b0b")).Bold(true)
)

type Tmodel struct {
	pdfs      []string
	directory string
	cursor    int
	selected  map[int]struct{}
	logo      string
	Quit      bool
}

func MultiSelectModel(pdfs []string, directory string, logo string) Tmodel {
	return Tmodel{
		pdfs:      pdfs,
		directory: directory,
		selected:  make(map[int]struct{}),
		logo:      logo,
	}
}

func (m Tmodel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Tmodel) GetSelectedPDFs() []string {
	var selected []string

	for i := range m.selected {
		selected = append(selected, m.pdfs[i])
	}
	return selected
}

func (m Tmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quit = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.pdfs)-1 {
				m.cursor++
			}

		case "x", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "enter":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Tmodel) View() string {
	var b strings.Builder
	switch m.logo {
	case "merge":
		b.WriteString(defaultStyle.Render(logoMerge))
		if len(m.pdfs) <= 1 {
			fmt.Fprint(&b, "\n\n")
			b.WriteString(errorStyle.Render("Please provide more than one PDF files. Press 'esc' to exit."))
			return b.String()
		}
		fmt.Fprint(&b, "\n\n")
		b.WriteString(defaultStyle.Render("Which PDFs do you want to merge together?"))

	case "encrypt":
		b.WriteString(defaultStyle.Render(logoEncrypt))
		if len(m.pdfs) == 0 {
			fmt.Fprint(&b, "\n\n")
			b.WriteString(errorStyle.Render("No PDF files found. Press 'esc' to exit."))
			return b.String()
		}
		fmt.Fprint(&b, "\n\n")
		b.WriteString(defaultStyle.Render("Which PDFs do you want to Encrypt?"))
	}

	fmt.Fprint(&b, "\n")
	b.WriteString(focusedStyle.Render("Select with Space or 'x'"))
	fmt.Fprint(&b, "\n\n")
	b.WriteString(selectedStyle.Render("File location: ", m.directory))
	fmt.Fprint(&b, "\n\n")

	for i, choice := range m.pdfs {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = selectedStyle.Render("x")
			choice = selectedStyle.Render(choice)
		}

		choice = defaultStyle.Render(choice)

		b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice))
	}

	b.WriteString(focusedStyle.Render("\nPress enter to confirm, esc to quit."))
	fmt.Fprint(&b, "\n\n")

	return b.String()
}
