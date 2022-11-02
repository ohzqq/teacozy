package lists

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Widget interface {
	Focus() tea.Cmd
	Focused() bool
	Blur()
	Label() string
	Toggle() key.Binding
	Update(*List, tea.Msg) tea.Cmd
	View() string
	SetContent(string)
}
