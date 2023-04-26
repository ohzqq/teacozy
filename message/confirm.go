package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Confirm struct {
	Question string
	Func     func(bool) tea.Cmd
}

func GetConfirmation(q string, fn func(bool) tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return Confirm{
			Question: q,
			Func:     fn,
		}
	}
}
