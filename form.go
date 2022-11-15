package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	main state = iota
	form
	view
	edit
)

type state int

func (s state) String() string {
	switch s {
	case form:
		return "form"
	case view:
		return "view"
	case edit:
		return "edit"
	default:
		return "main"
	}
}

type Form struct {
	Model  *List
	Input  textarea.Model
	view   viewport.Model
	Fields *Fields
	Hash   map[string]string
	state  state
}

func NewForm(data FormData) *Form {
	fields := NewFields().SetData(data)
	m := Form{
		Fields: fields,
		view:   fields.Display(),
		Model:  fields.Edit(),
	}
	return &m
}

func (f *Form) Render() *Form {
	items := NewItems()
	for _, key := range f.Fields.Data.Keys() {
		field := f.Fields.Data.Get(key)
		item := NewItem().SetData(field)
		items.Add(item)
	}
	m := NewList("Edit...", items)
	m.InitList()
	f.Model = m
	return f
}

func (m *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		if m.Input.Focused() {
			if key.Matches(msg, Keys.SaveAndExit) {
				cur := m.Model.Model.SelectedItem()
				i := m.Model.Items.Get(cur)
				field := i.Item.(FieldData)
				val := m.Input.Value()
				field.Set(val)
				item := NewItem().SetData(field)
				m.Model.Items.Set(i.Index(), item)
				m.Input.Blur()
				m.Render()
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch m.state {
			case view:
				switch {
				case key.Matches(msg, Keys.SaveAndExit):
					cmds = append(cmds, SaveFormAsHashCmd())
				case key.Matches(msg, Keys.EditField):
					cmds = append(cmds, EditFormCmd())
				case key.Matches(msg, Keys.ExitScreen):
					m.state = form
				}
				m.view, cmd = m.view.Update(msg)
				cmds = append(cmds, cmd)
			case form:
				switch {
				case key.Matches(msg, Keys.SaveAndExit):
					m.state = view
				//case key.Matches(msg, Keys.EditField):
				//cur := m.Model.Model.SelectedItem()
				//field := m.Model.Items.Get(cur).Item.(*Item) //.(*Field)
				//cmds = append(cmds, EditFormItemCmd(field))
				case key.Matches(msg, Keys.ExitScreen):
					cmds = append(cmds, tea.Quit)
				}
				m.Model.Model, cmd = m.Model.Model.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case SaveFormAsHashMsg:
		m.Hash = make(map[string]string)
		for _, field := range m.Fields.All() {
			m.Hash[field.Key()] = field.Value()
		}
		cmds = append(cmds, tea.Quit)
	case EditFormMsg:
		m.Render()
		m.state = form
	case EditFormItemMsg:
		m.Input = textarea.New()
		m.Input.SetValue(msg.Value())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case UpdateFormContentMsg:
		field := m.Fields.Data.Get(msg.Key())
		field.Set(msg.Value())
	case tea.WindowSizeMsg:
		m.Model.Model.SetSize(msg.Width-2, msg.Height-2)
	case ReturnSelectionsMsg:
		var field FieldData
		if items := m.Model.Items.Selections(); len(items) > 0 {
			field = items[0].Item.(*Field)
		}
		cmds = append(cmds, EditFormItemCmd(field.(*Item)))
	}

	return m, tea.Batch(cmds...)
}

func (m *Form) View() string {
	var (
		sections    []string
		availHeight = TermHeight()
	)
	var field string
	if m.Input.Focused() {
		field = m.Input.View()
		availHeight -= lipgloss.Height(field)
	}

	switch m.state {
	case view:
		v := m.view.View()
		sections = append(sections, v)
	case form:
		m.Model.Model.SetSize(m.Model.Width, availHeight)
		v := m.Model.View()
		sections = append(sections, v)
	}

	if m.Input.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
