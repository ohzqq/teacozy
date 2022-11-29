package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	*info.Info
	key.KeyMap
	Model viewport.Model
}

func NewHelp() Help {
	i := info.New()
	i.Toggle = key.NewKey("?", "help")
	km := KeyMap()
	i.NewSection().SetTitle("Help").SetFields(km)
	i.NewSection().SetTitle("Navigation").SetFields(ListKeyMap())
	//km := key.NewKeyMap()
	//km.Add(i.Toggle)
	return Help{
		Info:   i,
		KeyMap: km,
	}
}

func (h Help) View() string {
	h.Model.SetContent(h.Info.Render())
	return h.Model.View()
}

func (m Help) Update(msg tea.Msg) (Help, tea.Cmd) {
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
			m.Info.Hide()
			//cmds = append(cmds, ToggleHelpCmd())
		}
	}

	//m.Model = m.Info.Model
	m.Model, cmd = m.Info.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func GoToHelp(m tea.Model) tea.Cmd {
	if ui, ok := m.(*TUI); ok {
		return ui.ShowHelp()
	}
	return nil
}

func ListKeyMap() key.KeyMap {
	lk := list.ListKeyMap()
	km := key.NewKeyMap()
	km.AddBind(lk.CursorUp)
	km.AddBind(lk.CursorDown)
	return km
}
