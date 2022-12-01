package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
)

func (m *Tui) updateInfo(msg tea.Msg, nfo *info.Info) tea.Cmd {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		model tea.Model
	)

	switch msg := msg.(type) {
	case info.HideInfoMsg:
		m.showInfo = false
		m.showFullHelp = false
		m.state = mainModel
	case info.UpdateContentMsg:
		m.view.SetContent(msg.Content)
	case tea.KeyMsg:
		//cmds = append(cmds, list.UpdateStatusCmd(msg.String()))
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
	}

	model, cmd = nfo.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case infoModel:
		m.Info = model.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(m.Info.Render()))
	case helpModel:
		m.Help.Info = model.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(m.Help.Render()))
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

type Info struct {
	*info.Info
	showHelp   bool
	fullScreen bool
}

func NewInfo() *Info {
	i := info.New()
	return &Info{
		Info: i,
	}
}
