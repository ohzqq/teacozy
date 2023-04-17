package edit

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	KeyMap keys.KeyMap
	Prefix string
	help   keys.KeyMap
}

type Props struct {
	Value    string
	Save     func(string)
	ShowHelp func([]map[string]string)
}

type StartEditingMsg struct{}
type SaveEditMsg struct{}
type ConfirmEditMsg struct{}

func New() *Component {
	c := &Component{
		input:  textarea.New(),
		KeyMap: DefaultKeyMap(),
		Prefix: "> ",
	}
	c.input.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(color.Cyan())

	c.help = append(c.help, c.KeyMap...)
	c.help = append(c.help, keys.TextArea()...)

	return c
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prefix
	c.input.SetValue(props.Value)
	c.input.ShowLineNumbers = false
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ConfirmEditMsg:
		if c.Props().Value != c.input.Value() {
			c.Props().Save(c.input.Value())
			c.input.Reset()
			return confirm.GetConfirmation("Save edit?", SaveEdit)
		}
		return keys.ReturnToList
	case keys.ShowHelpMsg:
		c.Props().ShowHelp(c.help.Map())
		cmds = append(cmds, keys.ChangeRoute("help"))
	case tea.KeyMsg:
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	c.input, cmd = c.input.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	c.input.SetWidth(w)
	c.input.SetHeight(c.input.LineInfo().Height + 1)
	return c.input.View()
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
