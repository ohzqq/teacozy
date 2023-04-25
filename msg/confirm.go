package msg

import tea "github.com/charmbracelet/bubbletea"

type ConfirmFunc func(bool) tea.Cmd

type Confirm struct {
	Question string
	Confirm  ConfirmFunc
}

func GetConfirmation(q string, c Confirm) tea.Cmd {
	return func() tea.Msg {
		return Confirm{
			Question: q,
			Confirm:  c,
		}
	}
}
