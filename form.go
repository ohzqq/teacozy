package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	Frame
	*Items
	Model       list.Model
	Input       textarea.Model
	Info        *Info
	Fields      *Fields
	Confirm     *Menu
	Title       string
	FormChanged bool
	confirm     bool
	SaveForm    SaveForm
	Hash        map[string]string
	Style       list.Styles
}

func NewForm(fields *Fields) *Form {
	m := Form{
		SaveForm: SaveFormAsHash,
		Frame:    DefaultFrameStyle(),
		Items:    NewItems().SetShowKeys(),
		Fields:   fields,
		Confirm:  ConfirmMenu(),
	}
	m.Frame.MinHeight = 10

	for _, field := range fields.All() {
		i := NewItem().SetData(field)
		m.Add(i)
	}

	m.Model = NewListModel(m.Frame.Width(), m.Frame.Height(), m.Items)
	return &m
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
			if Keys.SaveAndExit.Matches(msg) {
				cur := m.Model.SelectedItem().(*Item)
				field := cur.Data
				val := m.Input.Value()
				if original := field.Value(); original != val {
					field.Set(val)
					item := NewItem().SetData(field)
					cmds = append(cmds, ItemChangedCmd(item))
				}
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case m.confirm:
				var mod tea.Model
				mod, cmd = m.Confirm.Update(msg)
				m.Confirm = mod.(*Menu)
			case m.FormChanged:
				switch {
				case Keys.SaveAndExit.Matches(msg):
					m.Confirm.AddContent(m.Fields.String())
					m.confirm = true
				}
			}
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	case UpdateStatusMsg:
		cmds = append(cmds, m.Model.NewStatusMessage(msg.Msg))
	case tea.WindowSizeMsg:
		m.Frame.SetSize(msg.Width-1, msg.Height-2)
		m.Model.SetSize(msg.Width-1, msg.Height-2)
	case EditFormItemMsg:
		cur := m.Model.SelectedItem().(*Item)
		m.Input = textarea.New()
		m.Input.SetValue(cur.Value())
		m.Input.ShowLineNumbers = false
		m.Input.Focus()
	case FormChangedMsg:
		m.FormChanged = true
	case ItemChangedMsg:
		idx := m.Model.Index()
		msg.Item.Changed = true
		m.Model.SetItem(idx, msg.Item)
		cmds = append(cmds, FormChangedCmd())
	case SetItemsMsg:
		m.Model.SetItems(msg.Items)
	case HideMenuMsg:
		println("hide")
		m.confirm = false
	case ConfirmMenuMsg:
		if msg {
			m.Confirm.Model.SetContent("confirmed")
			cmds = append(cmds, SaveFormCmd(m.SaveForm))
		}
	//case SaveAndExitFormMsg:
	//cmds = append(cmds, msg.Save(m))
	case SaveAndExitFormMsg:
		cmds = append(cmds, msg.Save(m))
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m *Form) Init() tea.Cmd {
	return nil
}

func (m Form) View() string {
	var (
		sections    []string
		availHeight = m.Frame.Height()
		field       string
	)

	if m.confirm {
		//info := m.Info.View()
		info := m.Confirm.View()
		//availHeight -= m.Info.Model.Height
		sections = append(sections, info)
	} else {
		if m.Input.Focused() {
			iHeight := availHeight / 3
			m.Input.SetHeight(iHeight)
			field = m.Input.View()
			availHeight -= iHeight
		}

		m.Frame.SetSize(m.Frame.Width(), availHeight)
		m.Model.SetSize(m.Frame.Width(), availHeight)
		content := m.Model.View()
		sections = append(sections, content)

		if m.Input.Focused() {
			sections = append(sections, field)
		}
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
