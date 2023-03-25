package react

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/style"
)

type FilterState int

// Possible filter states.
const (
	Unfiltered FilterState = iota // no filter set
	Filtering                     // user is actively setting a filter
)

type ListComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ListProps]

	*list.Model
}

type ListProps struct {
	Choices []string
}

func NewList() *ListComponent {
	return &ListComponent{
		Model: list.New([]string{}),
	}
}

func NewListProps(c []string) ListProps {
	return ListProps{Choices: c}
}

func (m *ListComponent) Init(props ListProps) tea.Cmd {
	m.Choices = props.Choices
	m.Items = list.ChoicesToMatch(props.Choices)
	m.Matches = m.Items
	m.Selected = make(map[int]struct{})

	m.Input = textinput.New()
	m.Input.Prompt = m.Prompt
	m.Input.PromptStyle = m.Style.Prompt
	m.Input.Placeholder = m.Placeholder
	m.Input.Width = 40

	v := viewport.New(40, 5)
	m.Viewport = &v

	m.Paginator = paginator.New()
	m.Paginator.SetTotalPages((len(m.Items) + 5 - 1) / 5)
	m.Paginator.PerPage = 5
	m.Paginator.Type = paginator.Dots
	m.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	m.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)

	m.UpdateProps(props)
	return nil
}

func (c *ListComponent) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	var cmd tea.Cmd
	_, cmd = c.Model.Update(msg)

	return cmd
}

func (m *ListComponent) Render(w, h int) string {
	return m.Model.View()
}
