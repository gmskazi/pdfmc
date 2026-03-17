package textInputs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gmskazi/pdfmc/cmd/styles"
)

var noStyle = lipgloss.NewStyle()

type Tmodel struct {
	focusIndex int
	inputs     []textinput.Model
	Quit       bool
}

func TextinputModel() Tmodel {
	m := Tmodel{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = styles.DefaultStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
			t.Focus()

		case 1:
			t.Placeholder = "Confirm Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.inputs[i] = t
	}

	return m
}

func (m Tmodel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tea.ClearScreen,
	)
}

func (m Tmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quit = true
			return m, tea.Quit

			// Set focus to next input
		case "tab", "shift-tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				if m.checkPasswords(m.inputs[0].Value(), m.inputs[1].Value()) {
					return m, tea.Quit
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift-tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// set focus state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.SelectedStyle
					m.inputs[i].TextStyle = styles.SelectedStyle
					continue
				}

				if m.checkPasswords(m.inputs[0].Value(), m.inputs[1].Value()) {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = styles.SelectedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = styles.ErrorStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m Tmodel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Tmodel) View() string {
	var b strings.Builder
	b.WriteString(styles.DefaultStyle.Render("Input the password to encrypt the PDFs."))
	fmt.Fprint(&b, "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(styles.FocusedStyle.Render("Press Enter on 'Submit' to continue"))
	fmt.Fprint(&b, "\n")

	b.WriteString(styles.FocusedStyle.Render("To Exit press 'ctrl+c' or 'esc'"))
	fmt.Fprint(&b, "\n")

	return b.String()
}

func (m Tmodel) GetPassword() string {
	return m.inputs[0].Value()
}

func (m Tmodel) checkPasswords(password, passwordConfirmation string) bool {
	return password == passwordConfirmation
}

func TextinputInteractive() (password string, quit bool, err error) {
	p := tea.NewProgram(TextinputModel())
	result, err := p.Run()
	if err != nil {
		return "", false, err
	}

	tmodel := result.(Tmodel)
	if tmodel.Quit {
		return "", true, fmt.Errorf("operation canceled")
	}

	pword := tmodel.GetPassword()
	return pword, false, nil
}
