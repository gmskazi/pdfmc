package multiSelect

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	logoMerge = `
 __  __                  
|  \/  |___ _ _ __ _ ___ 
| |\/| / -_) '_/ _` + "`" + ` / -_)
|_|  |_\___|_| \__, \___|
               |___/     
`

	logoEncrypt = `
 ___                       _   
| __|_ _  __ _ _ _  _ _ __| |_ 
| _|| ' \/ _| '_| || | '_ \  _|
|___|_||_\__|_|  \_, | .__/\__|
                 |__/|_|       
`

	logoDecrypt = `
 ___                       _   
|   \ ___ __ _ _ _  _ _ __| |_ 
| |) / -_) _| '_| || | '_ \  _|
|___/\___\__|_|  \_, | .__/\__|
                 |__/|_|       
`
	merge   = "merge"
	encrypt = "encrypt"
	decrypt = "decrypt"
)

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
	autoQuit  bool
	ErrMsg    string
}

func MultiSelectModel(pdfs []string, directory string, logo string) Tmodel {
	return Tmodel{
		pdfs:      pdfs,
		directory: directory,
		selected:  make(map[int]struct{}),
		logo:      logo,
	}
}

type autoQuitMsg struct{}

func (m Tmodel) Init() tea.Cmd {
	// Set error and autoQuit if conditions aren't met
	if m.logo == merge && len(m.pdfs) <= 1 {
		return func() tea.Msg {
			return autoQuitMsg{}
		}
	} else if m.logo == encrypt && len(m.pdfs) == 0 {
		return func() tea.Msg {
			return autoQuitMsg{}
		}
	} else if m.logo == decrypt && len(m.pdfs) == 0 {
		return func() tea.Msg {
			return autoQuitMsg{}
		}
	}
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
	case autoQuitMsg:
		m.autoQuit = true
		if m.logo == merge && len(m.pdfs) <= 1 {
			m.ErrMsg = "Error: Need at least 2 PDFs to merge"
		} else {
			m.ErrMsg = "Error: No PDFs found to encrypt"
		}
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quit = true
			return m, tea.Quit

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "j", "down":
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
	case merge:
		b.WriteString(defaultStyle.Render(logoMerge))
		fmt.Fprint(&b, "\n\n")
	case encrypt:
		b.WriteString(defaultStyle.Render(logoEncrypt))
		fmt.Fprint(&b, "\n\n")
	case decrypt:
		b.WriteString(defaultStyle.Render(logoDecrypt))
		fmt.Fprint(&b, "\n\n")
	}

	if m.ErrMsg != "" {
		b.WriteString(errorStyle.Render(m.ErrMsg))
		return b.String()
	}

	switch m.logo {
	case merge:
		b.WriteString(defaultStyle.Render("Which PDFs do you want to merge together?"))

	case encrypt:
		b.WriteString(defaultStyle.Render("Which PDFs do you want to Encrypt?"))

	case decrypt:
		b.WriteString(defaultStyle.Render("Which PDFs do you want to Decrypt?"))
	}

	fmt.Fprint(&b, "\n")
	b.WriteString(focusedStyle.Render("Select with Space or 'x', navigate with up/down or j/k"))
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

func MultiSelectInteractive(pdfs []string, dir string, logo string) (selectedPdfs []string, quit bool, err error) {
	p := tea.NewProgram(MultiSelectModel(pdfs, dir, logo))
	result, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	model := result.(Tmodel)
	if model.autoQuit {
		return nil, true, fmt.Errorf("%s", model.ErrMsg)
	}

	if model.Quit {
		return nil, true, fmt.Errorf("operation canceled")
	}

	selectedPdfs = model.GetSelectedPDFs()

	return selectedPdfs, false, nil
}
