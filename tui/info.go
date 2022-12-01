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

func (ui *Tui) updateInfo(msg tea.Msg, m *info.Info) tea.Cmd {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		model tea.Model
	)

	model, cmd = m.Update(msg)
	cmds = append(cmds, cmd)

	switch ui.state {
	case infoModel:
		ui.Info = model.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(ui.Info.Render()))
	case helpModel:
		ui.Help.Info = model.(*info.Info)
		cmds = append(cmds, info.UpdateContentCmd(ui.Help.Render()))
	}

	return tea.Batch(cmds...)
}

func (ui *Tui) viewInfo() string {
	var (
		widgetWidth  = ui.Style.Widget.Width()
		widgetHeight = ui.Style.Widget.Height()
	)
	ui.view = viewport.New(widgetWidth, widgetHeight)
	switch ui.state {
	case infoModel:
	case helpModel:
	}
	return ui.Info.View()
}
