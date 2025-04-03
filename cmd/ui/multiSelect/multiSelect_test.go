package multiSelect

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestMultiSelectModel(t *testing.T) {
	tests := []struct {
		name          string
		pdfs          []string
		directory     string
		logo          string
		expectedModel Tmodel
	}{
		{
			name:      "Merge model initialization",
			pdfs:      []string{"a.pdf", "b.pdf"},
			directory: "/test",
			logo:      "merge",
			expectedModel: Tmodel{
				pdfs:      []string{"a.pdf", "b.pdf"},
				directory: "/test",
				selected:  make(map[int]struct{}),
				logo:      "merge",
			},
		},
		{
			name:      "Encrypt model initialization",
			pdfs:      []string{"secret.pdf"},
			directory: "/secure",
			logo:      "encrypt",
			expectedModel: Tmodel{
				pdfs:      []string{"secret.pdf"},
				directory: "/secure",
				selected:  make(map[int]struct{}),
				logo:      "encrypt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := MultiSelectModel(tt.pdfs, tt.directory, tt.logo)
			assert.Equal(t, tt.expectedModel.pdfs, model.pdfs)
			assert.Equal(t, tt.expectedModel.directory, model.directory)
			assert.Equal(t, tt.expectedModel.logo, model.logo)
			assert.Equal(t, 0, model.cursor)
			assert.Empty(t, model.selected)
		})
	}
}

func TestGetSelectedPDFs(t *testing.T) {
	tests := []struct {
		name     string
		selected map[int]struct{}
		pdfs     []string
		expected []string
	}{
		{
			name:     "No selections",
			selected: map[int]struct{}{},
			pdfs:     []string{"a.pdf", "b.pdf"},
			expected: nil,
		},
		{
			name:     "Single selection",
			selected: map[int]struct{}{0: {}},
			pdfs:     []string{"a.pdf", "b.pdf"},
			expected: []string{"a.pdf"},
		},
		{
			name:     "Multiple selections",
			selected: map[int]struct{}{0: {}, 2: {}},
			pdfs:     []string{"a.pdf", "b.pdf", "c.pdf"},
			expected: []string{"a.pdf", "c.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Tmodel{
				pdfs:     tt.pdfs,
				selected: tt.selected,
			}
			assert.Equal(t, tt.expected, model.GetSelectedPDFs())
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		initial    Tmodel
		msg        tea.Msg
		expected   Tmodel
		shouldQuit bool
	}{
		{
			name: "Move cursor up",
			initial: Tmodel{
				pdfs:   []string{"a.pdf", "b.pdf"},
				cursor: 1,
			},
			msg: tea.KeyMsg{Type: tea.KeyUp},
			expected: Tmodel{
				pdfs:   []string{"a.pdf", "b.pdf"},
				cursor: 0,
			},
		},
		{
			name: "Move cursor down",
			initial: Tmodel{
				pdfs:   []string{"a.pdf", "b.pdf"},
				cursor: 0,
			},
			msg: tea.KeyMsg{Type: tea.KeyDown},
			expected: Tmodel{
				pdfs:   []string{"a.pdf", "b.pdf"},
				cursor: 1,
			},
		},
		{
			name: "Select item with space",
			initial: Tmodel{
				pdfs:     []string{"a.pdf", "b.pdf"},
				selected: make(map[int]struct{}),
			},
			msg: tea.KeyMsg{Type: tea.KeySpace},
			expected: Tmodel{
				pdfs:     []string{"a.pdf", "b.pdf"},
				selected: map[int]struct{}{0: {}},
			},
		},
		{
			name: "Deselect item with space",
			initial: Tmodel{
				pdfs:     []string{"a.pdf", "b.pdf"},
				selected: map[int]struct{}{0: {}},
			},
			msg: tea.KeyMsg{Type: tea.KeySpace},
			expected: Tmodel{
				pdfs:     []string{"a.pdf", "b.pdf"},
				selected: map[int]struct{}{},
			},
		},
		{
			name: "Quit with escape",
			initial: Tmodel{
				pdfs: []string{"a.pdf", "b.pdf"},
			},
			msg: tea.KeyMsg{Type: tea.KeyEsc},
			expected: Tmodel{
				pdfs: []string{"a.pdf", "b.pdf"},
				Quit: true,
			},
			shouldQuit: true,
		},
		{
			name: "Confirm with enter",
			initial: Tmodel{
				pdfs: []string{"a.pdf", "b.pdf"},
			},
			msg: tea.KeyMsg{Type: tea.KeyEnter},
			expected: Tmodel{
				pdfs: []string{"a.pdf", "b.pdf"},
			},
			shouldQuit: true,
		},
		{
			name: "Auto quit for merge with insufficient PDFs",
			initial: Tmodel{
				pdfs: []string{"a.pdf"},
				logo: "merge",
			},
			msg: autoQuitMsg{},
			expected: Tmodel{
				pdfs:     []string{"a.pdf"},
				logo:     "merge",
				autoQuit: true,
				ErrMsg:   "Error: Need at least 2 PDFs to merge",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, cmd := tt.initial.Update(tt.msg)

			// Check if we got a quit command
			if tt.shouldQuit {
				if cmd != nil {
					_, ok := cmd().(tea.QuitMsg)
					assert.True(t, ok, "Expected quit command")
				} else {
					assert.True(t, model.(Tmodel).Quit, "Expected Quit flag to be set")
				}
			}

			assert.Equal(t, tt.expected, model.(Tmodel))
		})
	}
}
