package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	urkey "github.com/ohzqq/teacozy/key"
)

type Delegate struct {
	list.DefaultDelegate
}

func NewDelegate() Delegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.UpdateFunc = func(msg tea.Msg, model *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, urkey.EditField):
				return model.NewStatusMessage("edit field")
				//cmds = append(cmds, EditContentCmd(curItem))
			}
		}
		return nil
	}
	return Delegate{DefaultDelegate: d}
}
