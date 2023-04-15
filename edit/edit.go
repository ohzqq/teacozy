package edit

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/app/confirm"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	Placeholder string
	Prompt      string
	Style       style.List
	KeyMap      keys.KeyMap
}

type Props struct {
	Value string
	Save  func(string)
}

type StartEditingMsg struct{}
type SaveEditMsg struct{}
type ConfirmEditMsg struct{}

func New() *Component {
	tm := &Component{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
		input:  textarea.New(),
		KeyMap: DefaultKeyMap(),
	}
	return tm
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prompt
	c.input.Placeholder = c.Placeholder
	c.input.SetValue(props.Value)
	c.input.ShowLineNumbers = false
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SaveEditMsg:
		if c.Props().Value != c.input.Value() {
			c.Props().Save(c.input.Value())
			c.input.Reset()
			return confirm.GetConfirmation("Save edit?", SaveEdit)
		}
		return keys.ReturnToList
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
	c.input.SetHeight(c.input.LineInfo().Height)
	return c.input.View()
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Esc(),
		keys.Quit(),
		keys.Save().Cmd(Save),
	}
	return km
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
