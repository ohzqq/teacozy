package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
)

type Info struct {
	*info.Info
	Model    viewport.Model
	Help     Help
	Toggle   *key.Key
	showHelp bool
}

func NewInfo() *Info {
	i := info.New()
	return &Info{
		Info: i,
	}
}

func (m *Info) Update(msg tea.Msg) (*Info, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, key.HelpKey):
			cmds = append(cmds, ToggleHelpCmd())
			//cmds = append(cmds, info.HideInfoCmd())
		}
	case tea.WindowSizeMsg:
		m.Model = viewport.New(msg.Width-2, msg.Height-2)
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Info) View() string {
	var content string
	if m.showHelp {
		content = m.Help.Render()
	}
	m.Model.SetContent(content)
	return m.Model.View()
}
