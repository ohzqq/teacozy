package info

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
)

var fieldStyle = Style{
	Key:   lipgloss.NewStyle().Foreground(style.DefaultColors().Blue),
	Value: lipgloss.NewStyle().Foreground(style.DefaultColors().DefaultFg),
}

type Style struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

func SetKeyStyle(s lipgloss.Style) {
	fieldStyle.Key = s
}

func SetValueStyle(s lipgloss.Style) {
	fieldStyle.Value = s
}

func (i *Fields) NoKeys() *Fields {
	i.HideKeys = true
	return i
}

func (m *Fields) Update(msg tea.Msg) (*Fields, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (i Fields) String() string {
	var info []string
	for _, key := range i.Keys() {
		var line []string
		field := i.Get(key)
		if !i.HideKeys {
			k := i.Style.Key.Render(field.Key())
			line = append(line, k, ": ")
		}

		v := i.Style.Value.Render(field.Value())
		line = append(line, v)

		l := strings.Join(line, "")
		info = append(info, l)
	}
	return strings.Join(info, "\n")
}

func (m *Fields) View() string {
	m.Model.SetContent(m.String())
	return m.Model.View()
}

func (i *Fields) Init() tea.Cmd {
	return nil
}
