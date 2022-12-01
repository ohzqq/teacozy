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

func updateInfo(msg tea.Msg, m Tui) (viewport.Model, tea.Cmd) {
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
	m.Info.Model.SetContent(m.Info.Render())
	i, cmd = m.Info.Update(msg)
	m.Info = i.(*info.Info)
	cmds = append(cmds, cmd)

	//m.view = m.Info.Model
	//m.view.SetContent(m.Info.Render())
	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return m.Info.Model, tea.Batch(cmds...)
}

func (m *Info) Update(msg tea.Msg, ui Tui) (*Info, tea.Cmd) {
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

	var model tea.Model
	switch m.state {
	case infoModel:
		m.Model.SetContent(m.Info.Render())
		model, cmd = ui.Info.Update(msg)
		ui.Info = model.(*info.Info)
	case helpModel:
		m.Model.SetContent(m.Help.Render())
		model, cmd = ui.Help.Info.Update(msg)
		ui.Help.Info = model.(*info.Info)
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
