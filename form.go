package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	form state = iota
	view
	edit
)

type Form struct {
	Model  *List
	Input  textarea.Model
	Fields *Fields
	Hash   map[string]string
	state  state
}

func New(data FormData) *Form {
	fields := NewFields().SetData(data)
	m := Form{
		Fields: fields,
	}
	m.Fields.Render()
	m.Render()
	return &m
}

func (f *Form) Render() *Form {
	items := NewItems()
	for _, key := range f.Fields.Data.Keys() {
		field := f.Fields.Data.Get(key)
		item := NewItem(field)
		items.Add(item)
	}
	m := NewList()
	m.Title = "Edit..."
	m.SetItems(items).InitList()
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
				cur := m.Model.List.SelectedItem()
				i := m.Model.Items.Get(cur)
				field := i.Data.(Field)
				val := m.Input.Value()
				field.Set(val)
				m.Model.Items.Set(i.Index(), NewItem(field))
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
					cmds = append(cmds, SaveAsHashCmd())
				case key.Matches(msg, Keys.EditField):
					cmds = append(cmds, EditInfoCmd())
				case key.Matches(msg, Keys.ExitScreen):
					m.state = form
				}
				m.Fields, cmd = m.Fields.Update(msg)
				cmds = append(cmds, cmd)
			case form:
				switch {
				case key.Matches(msg, Keys.SaveAndExit):
					m.state = view
				case key.Matches(msg, Keys.EditField):
					cur := m.Model.List.SelectedItem()
					field := m.Model.Items.Get(cur).Data.(Field)
					cmds = append(cmds, EditItemCmd(field))
				case key.Matches(msg, Keys.ExitScreen):
					cmds = append(cmds, tea.Quit)
				}
				m.Model.List, cmd = m.Model.List.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case SaveAsHashMsg:
		m.Hash = make(map[string]string)
		for _, field := range m.Fields.AllFields() {
			m.Hash[field.Key()] = field.Value()
		}
		cmds = append(cmds, tea.Quit)
	case EditInfoMsg:
		m.Render()
		m.state = form
	case EditItemMsg:
		m.Input = textarea.New()
		m.Input.SetValue(msg.Value())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case UpdateContentMsg:
		field := m.Fields.Data.Get(msg.Key())
		field.Set(msg.Value())
	case tea.WindowSizeMsg:
		m.Model.List.SetSize(msg.Width-2, msg.Height-2)
	case ReturnSelectionsMsg:
		var field Field
		if items := m.Model.Items.Selections(); len(items) > 0 {
			field = items[0].Data.(Field)
		}
		cmds = append(cmds, EditItemCmd(field))
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
		v := m.Fields.View()
		sections = append(sections, v)
	case form:
		m.Model.List.SetSize(m.Model.Width, availHeight)
		v := m.Model.View()
		sections = append(sections, v)
	}

	if m.Input.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
