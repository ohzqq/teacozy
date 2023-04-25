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
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	KeyMap      keys.KeyMap
	Prompt      string
	help        keys.KeyMap
	current     int
	originalVal string
	newVal      string
}

type Props struct {
	ShowHelp func([]map[string]string)
	teacozy.Props
}

func New() *Component {
	c := &Component{
		input:  textarea.New(),
		Prompt: "> ",
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
	c.input.Blur()
	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case keys.ConfirmEditMsg:
		//c.Props().Save(c.input.Value())
		//c.input.Reset()
		cmds = append(cmds, confirm.GetConfirmation("Save edit?", c.SaveEdit(msg.Value)))
		//cmds = append(cmds, keys.UpdateStatus("confirm"))
		//cmds = append(cmds, keys.ReturnToList)
		return tea.Batch(cmds...)
	case keys.StopEditingMsg:
		//c.input.Reset()
		c.input.Blur()
		c.Props().SetKeyMap(frame.DefaultKeyMap())
		c.newVal = c.input.Value()
		if c.newVal != c.originalVal {
			return keys.ConfirmEdit(c.input.Value())
		}
		return nil
	case keys.EditItemMsg:
		c.current = msg.Index
		c.originalVal = c.Props().Items.String(c.current)
		c.Props().SetKeyMap(pagy.DefaultKeyMap())
		c.input.SetValue(c.originalVal)
		return c.input.Focus()
	//case keys.SaveChangesMsg:
	//return keys.UpdateStatus(strconv.Itoa(msg.Index))
	//return keys.UpdateStatus(c.newVal)
	//case keys.SaveEditMsg:
	//val := c.Props().Items.String(c.current)
	//if in := c.input.Value(); in != val {
	//c.Props().Items.Set(c.current, in)
	//}
	//return keys.StopEditing
	case tea.KeyMsg:
		for _, k := range c.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
		if c.input.Focused() {
			c.input, cmd = c.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (c *Component) SaveEdit(v string) func(bool) tea.Cmd {
	return func(save bool) tea.Cmd {
		if save {
			fn := func(idx int) tea.Cmd {
				c.Props().Items.Set(idx, v)
				return nil
			}
			return keys.UpdateItem(fn)
		}
		return keys.ReturnToList
	}
}

func (c *Component) Render(w, h int) string {
	props := c.Props().Props

	if c.input.Focused() {
		c.input.SetWidth(w)
		ih := c.input.LineInfo().Height + 1
		h = h - ih
		c.input.SetHeight(ih)
	}
	input := c.input.View()

	view := teacozy.Renderer(props, w, h)
	if !c.input.Focused() {
		return view
	}

	return lipgloss.JoinVertical(lipgloss.Left, view, input)
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["form"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			Props: a.ItemProps(),
		}
		//a.SetKeyMap(pagy.DefaultKeyMap())
		a.SetKeyMap(frame.DefaultKeyMap())
		return comp, comp.Init(p)
	}
}

func (c Component) Name() string {
	return "form"
}

func DefaultKeyMap() keys.KeyMap {
	k := []*keys.Binding{
		keys.Esc().Cmd(keys.StopEditing),
		keys.Quit(),
		keys.Help(),
		keys.Edit(),
		keys.Save().
			Cmd(keys.StopEditing),
	}
	km := keys.NewKeyMap(k...)
	return km
}

func DefaultStyle() textarea.Style {
	return textarea.Style{
		Base:       lipgloss.NewStyle(),
		CursorLine: lipgloss.NewStyle().Background(color.Grey()),
	}
}
