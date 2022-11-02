package prompt

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/list"
)

type Prompt struct {
	*list.Model
	data []string
}

func New(title string, items []string, multi bool) *Prompt {
	l := list.New(title)
	l.IsPrompt = true

	if multi {
		l.SetMulti()
	}

	for _, key := range items {
		i := list.NewDefaultItem(key, key)
		l.AppendItem(i)
	}

	l.List = l.BuildModel()

	return &Prompt{
		List: l,
		data: items,
	}
}

func (m *Prompt) Choose() []string {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	var choices []string
	for _, sel := range m.Selections {
		choices = append(choices, sel.Title())
	}

	return choices
}
