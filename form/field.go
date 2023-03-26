package form

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type Field struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FieldProps]

	input   textinput.Model
	padding int
}

type FieldProps struct {
	idx      int
	key      string
	val      string
	name     string
	input    textinput.Model
	SetValue func(int, string)
}

//func (f *Field) Init(props FieldProps) tea.Cmd {
//  f.UpdateProps(props)
//  return f.input.Focus()
//}

func (c *Field) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.Type == tea.KeyEnter {
			// Lifted state power! Woohooo
			c.Props().SetValue(c.Props().idx, c.input.Value())

			// Navigate to displayname, please
			reactea.SetCurrentRoute("default")
			return nil
		}
	}
	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return cmd
}

func (f *Field) Render(w, h int) string {
	var s strings.Builder
	s.WriteString(f.Props().key)
	s.WriteString(": ")
	s.WriteString(f.input.View())
	s.WriteString("\n")
	return lipgloss.NewStyle().PaddingTop(f.padding).Render(s.String())
}

func NewForm() *Form {
	return &Form{}
}

func (f *Form) Init(props FormProps) tea.Cmd {
	f.UpdateProps(props)
	return nil
}

func (c *Form) Update(msg tea.Msg) tea.Cmd {
	cur := c.Props().Fields[1]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.Type == tea.KeyEnter {
			if cur.input.Focused() {
				c.Props().SetValue(cur.idx, cur.input.Value())
				cur.input.Blur()
			} else {
				return cur.input.Focus()
			}
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	cur.input, cmd = cur.input.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = c.view.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Form) Render(width, height int) string {
	view := viewport.New(width, height)
	var s strings.Builder
	for _, field := range c.Props().Fields {
		//for key, val := range field {
		s.WriteString("[")
		s.WriteString(field.key)
		s.WriteString("]")
		if field.input.Focused() {
			s.WriteString(field.input.View())
		} else {
			s.WriteString(field.input.Value())
		}
		s.WriteString("\n")
		//}
	}
	c.view = &view
	c.view.SetContent(s.String())
	return c.view.View()
}
