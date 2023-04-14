package edit

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
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

type StopEditingMsg struct{}
type StartEditingMsg struct{}
type SaveEditMsg struct{}

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
		c.Props().Save(c.input.Value())
		c.input.Reset()
		cmds = append(cmds, StopEditing)
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
	//l := lipgloss.NewStyle().Width(w).Render(c.Props().Value)
	//c.input.SetHeight(lipgloss.Height(l))
	c.input.SetWidth(w)
	c.input.SetHeight(c.input.LineInfo().Height)
	//fmt.Println(strconv.Itoa(len(c.Props().Value)))
	//fmt.Println(strconv.Itoa(c.input.Length()))
	return c.input.View()
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Esc().Cmd(StopEditing),
		keys.Quit(),
		keys.Save().Cmd(Save),
	}
	return km
}

func Save() tea.Msg {
	return SaveEditMsg{}
}

func StartEditing() tea.Msg {
	return StartEditingMsg{}
}

func StopEditing() tea.Msg {
	return StopEditingMsg{}
}
