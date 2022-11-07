package info

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list/item"
	"github.com/ohzqq/teacozy/prompt"
	"github.com/ohzqq/teacozy/util"
)

type Form struct {
	Model  prompt.Model
	Input  textarea.Model
	Fields *Fields
	state  state
}

func (f *Form) Edit() *Form {
	items := item.NewItems()
	for _, key := range f.Fields.Keys() {
		field := f.Fields.Get(key)
		item := item.NewItem(field)
		items.Add(item)
	}
	m := prompt.New()
	m.Title = "Edit..."
	m.SetItems(items).MakeList()
	f.Model = m
	return f
}

//func (m *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if key.Matches(msg, urkey.SaveAndExit) {
				cur := m.Model.List.SelectedItem()
				i := m.Model.Items.Get(cur)
				field := i.Data.(Field)
				val := m.Input.Value()
				field.Set(val)
				m.Model.Items.Set(i.Index(), item.NewItem(field))
				m.Input.Blur()
				m.Edit()
				cmds = append(cmds, prompt.UpdateVisibleItemsCmd("visible"))
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch m.state {
			case view:
				switch {
				case key.Matches(msg, urkey.EditField):
					cmds = append(cmds, EditInfoCmd())
				}
			case form:
				switch {
				case key.Matches(msg, urkey.EditField):
					cur := m.Model.List.SelectedItem()
					field := m.Model.Items.Get(cur).Data.(Field)
					cmds = append(cmds, EditItemCmd(field))
				case key.Matches(msg, urkey.ExitScreen):
					m.state = view
				}
			}
		}
	case EditInfoMsg:
		m.Edit()
		m.state = form
	case EditItemMsg:
		m.Input = textarea.New()
		m.Input.SetValue(msg.Value())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case UpdateContentMsg:
		m.Fields.Set(msg.Key(), msg.Value())
	case tea.WindowSizeMsg:
		m.Model.List.SetSize(msg.Width-2, msg.Height-2)
	case prompt.ReturnSelectionsMsg:
		var field Field
		if items := m.Model.Items.Selections(); len(items) > 0 {
			field = items[0].Data.(Field)
		}
		cmds = append(cmds, EditItemCmd(field))
	case prompt.UpdateStatusMsg:
	}

	switch m.state {
	case view:
		m.Fields, cmd = m.Fields.Update(msg)
		cmds = append(cmds, cmd)
	case form:
		m.Model.List, cmd = m.Model.List.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Form) View() string {
	var (
		sections []string
		//availHeight = m.form.List.Height()
		availHeight = util.TermHeight()
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

func (f *Form) Init() tea.Cmd {
	return nil
}
