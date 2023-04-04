package confirm

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
)

type Confirm struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Style    lipgloss.Style
	question string
}

type Props struct {
	*props.Items
}

func New() *Confirm {
	return &Confirm{
		Style: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
	}
}

func (c Confirm) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(Props{Items: props})
	}
}

func (c *Confirm) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Confirm) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
	case message.ReturnSelectionsMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		case "y":
			cmds = append(cmds, ConfirmCmd(true))
			cmds = append(cmds, message.ChangeRoute(c.Props().PrevRoute))
		case "n":
			cmds = append(cmds, ConfirmCmd(false))
			cmds = append(cmds, message.ChangeRoute(c.Props().PrevRoute))
		}
	}
	return tea.Batch(cmds...)
}

func (c *Confirm) Render(w, h int) string {
	return fmt.Sprintf("%s\n", c.Style.Render(c.Props().Question+"(y/n)"))
}

type GetConfirmationMsg struct {
	Question string
}

func GetConfirmation(q string) tea.Cmd {
	return func() tea.Msg {
		return GetConfirmationMsg{Question: q}
	}
}

type ConfirmMsg struct {
	Confirmed bool
}

func ConfirmCmd(confirm bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{Confirmed: confirm}
	}
}
