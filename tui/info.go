package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
)

type infoState int

type Info struct {
	*info.Info
	Model    viewport.Model
	Help     Help
	Toggle   *key.Key
	showHelp bool
	state    state
}

func NewInfo() *Info {
	i := info.New()
	return &Info{
		Info:  i,
		state: infoModel,
	}
}

func updateInfo(msg tea.Msg, m *Tui) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		//cmds = append(cmds, list.UpdateStatusCmd(msg.String()))
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
	}

	var i tea.Model
	switch m.state {
	case helpModel:
		i, cmd = m.Help.Update(msg)
		cmds = append(cmds, cmd)
		m.Help.Info = i.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(m.Help.Render()))
	default:
		i, cmd = m.Info.Update(msg)
		cmds = append(cmds, cmd)
		m.Info = i.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(m.Info.Render()))
	}

	return tea.Batch(cmds...)
}
