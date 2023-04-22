package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	KeyMap  keys.KeyMap
	Prompt  string
	help    keys.KeyMap
	current int
}

type Props struct {
	teacozy.Props
	ShowHelp func([]map[string]string)
}

type StartEditingMsg struct{}
type SaveEditMsg struct{}
type ConfirmEditMsg struct{}

func New() *Component {
	c := &Component{
		input:  textarea.New(),
		Prompt: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
		KeyMap: DefaultKeyMap(),
	}

	return c
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prompt
	c.input.KeyMap = keys.TextAreaDefault()
	c.input.ShowLineNumbers = false
	c.input.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(color.Cyan())
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ConfirmEditMsg:
		//if c.Props().Value != c.input.Value() {
		//c.Props().Save(c.input.Value())
		//c.input.Reset()
		//return confirm.GetConfirmation("Save edit?", SaveEdit)
		//}
		return keys.ReturnToList
	case tea.KeyMsg:
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	if c.input.Focused() {
		c.input, cmd = c.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	c.input.SetWidth(w)
	ih := c.input.LineInfo().Height + 1
	c.input.SetHeight(ih)
	view := c.input.View()
	props := c.Props().Props
	props.SetCurrent = c.setCurrent
	return lipgloss.JoinVertical(lipgloss.Left, view, teacozy.Renderer(props, w, h-ih))
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["form"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			Props:       a.ItemProps(),
			ToggleItems: a.ToggleItems,
		}
		a.SetKeyMap(pagy.DefaultKeyMap())
		return comp, comp.Init(p)
	}
}

func (c *Component) setCurrent(i int) {
	c.current = i
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Esc(),
		keys.Quit(),
		keys.Save().Cmd(ConfirmEdit),
		keys.Help(),
	}
	return km
}

func DefaultStyle() textarea.Style {
	return textarea.Style{
		Base:       lipgloss.NewStyle(),
		CursorLine: lipgloss.NewStyle().Background(color.Grey()),
	}
}

func Save() tea.Msg {
	return SaveEditMsg{}
}

func SaveEdit(save bool) tea.Cmd {
	if save {
		return Save
	}
	return keys.ReturnToList
}

func ConfirmEdit() tea.Msg {
	return ConfirmEditMsg{}
}

func StartEditing() tea.Msg {
	return StartEditingMsg{}
}
