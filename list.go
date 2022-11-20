package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionFunc func(items ...*Item) tea.Cmd

type List struct {
	Model         list.Model
	Input         textarea.Model
	Title         string
	SelectionList bool
	Editable      bool
	FormChanged   bool
	SaveFormFunc  SaveFormFunc
	ActionFunc    ActionFunc
	Hash          map[string]string
	Style         list.Styles
	id            int
	Frame
	*Items
}

func NewList() *List {
	m := List{
		SaveFormFunc: SaveFormAsHashCmd,
		Frame:        DefaultFrameStyle(),
		Items:        NewItems(),
	}
	m.SetAction(PrintItems)
	m.Frame.MinHeight = 10
	return &m
}

func NewListModel(w, h int, items *Items) list.Model {
	l := list.New(items.Visible(), items, w, h)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Styles = ListStyles()
	return l
}

func (m *List) ChooseOne() *List {
	m.SetModel()
	return m
}

func (m *List) ChooseMany() *List {
	m.SetMultiSelect()
	m.SetModel()
	return m
}

func (m *List) Edit() *List {
	m.Editable = true
	m.SetShowKeys()
	m.SetModel()
	return m
}

func (m *List) SetModel() *List {
	if !m.Editable {
		m.Items.Process()
	}
	m.Model = NewListModel(m.Width(), m.Height(), m.Items)
	return m
}

func (m *List) SetAction(fn ActionFunc) *List {
	m.ActionFunc = fn
	return m
}

func (m *List) SetItems(items *Items) *List {
	m.Items = items
	return m
}

func (m *List) SetTitle(title string) *List {
	m.Model.Title = title
	return m
}

func (m *List) SetSize(w, h int) *List {
	m.Frame.SetSize(w, h)
	m.Model.SetSize(m.Width(), m.Height())
	return m
}

func (m List) Height() int {
	return m.Frame.Height()
}

func (m List) Width() int {
	return m.Frame.Width()
}

func (m *List) SetMultiSelect() *List {
	m.Items.SetMultiSelect()
	return m
}

func (m *List) MultiSelect() bool {
	return m.Items.MultiSelect
}

func (m *List) ShowKeys() bool {
	return m.Items.ShowKeys
}

func (m *List) SetShowKeys() *List {
	m.Items.SetShowKeys()
	return m
}

func (m *List) SelectedItem() *Item {
	sel := m.Model.SelectedItem()
	cur := m.Items.Get(sel)
	return cur
}

func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
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
				cur := m.SelectedItem()
				field := cur.Data
				val := m.Input.Value()
				if original := field.Value(); original != val {
					field.Set(val)
					item := NewItem().SetData(field)
					cmds = append(cmds, ItemChangedCmd(item))
				}
				m.Input.Blur()
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case Keys.PrevScreen.Matches(msg):
				m.SelectionList = false
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			switch {
			case m.Editable:
				switch {
				case Keys.SaveAndExit.Matches(msg):
					cmds = append(cmds, FormChangedCmd())
				}
			case m.SelectionList:
				switch {
				case Keys.Enter.Matches(msg):
					cmds = append(cmds, ReturnSelectionsCmd())
				}
			case m.MultiSelect():
				switch {
				case Keys.Enter.Matches(msg):
					if m.SelectionList {
						cmds = append(cmds, ReturnSelectionsCmd())
					}
					m.SelectionList = true
					cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
				case Keys.DeselectAll.Matches(msg):
					m.Items.DeselectAllItems()
					cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
				case Keys.ToggleAllItems.Matches(msg):
					m.Items.ToggleAllSelectedItems()
					cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
				}
			default:
				switch {
				case Keys.Enter.Matches(msg):
					cur := m.Model.SelectedItem().(*Item)
					m.Items.ToggleSelectedItem(cur.Index())
					cmds = append(cmds, ReturnSelectionsCmd())
				}
			}
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)

		}
	case UpdateStatusMsg:
		cmds = append(cmds, m.Model.NewStatusMessage(msg.Msg))
	case tea.WindowSizeMsg:
		m.Model.SetSize(msg.Width-1, msg.Height-2)
	case UpdateVisibleItemsMsg:
		var items []list.Item
		switch string(msg) {
		case "selected":
			items = m.Selections()
		case "all":
			items = m.AllItems()
		default:
			items = m.Visible()
		}
		cmds = append(cmds, SetItemsCmd(items))
	case ToggleSelectedItemMsg:
		m.Items.ToggleSelectedItem(msg.Index())
	case ReturnSelectionsMsg:
		var items []*Item
		for _, sel := range m.Selections() {
			items = append(items, m.Get(sel))
		}
		//cmds = append(cmds, m.Model.ToggleSpinner())
		cmds = append(cmds, m.ActionFunc(items...))
	case EditFormItemMsg:
		if m.Editable {
			m.Input = textarea.New()
			m.Input.SetValue(msg.Value())
			m.Input.ShowLineNumbers = false
			m.Input.Focus()
		}
	case FormChangedMsg:
		m.FormChanged = true
	case ItemChangedMsg:
		msg.Item.Changed = true
		m.Items.Set(msg.Item.Index(), msg.Item)
		cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
	case SortItemsMsg:
		m.Items.SetItems(msg.Items...)
		m.Items.Process()
		cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
	case SetItemsMsg:
		m.Model.SetItems(msg.Items)
	case ToggleItemChildrenMsg:
		switch msg.ShowChildren {
		case true:
			m.Items.CloseItemList(msg.Index())
		default:
			m.Items.OpenItemList(msg.Index())
		}
		cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
	}

	return m, tea.Batch(cmds...)
}

func (m *List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	var (
		sections    []string
		availHeight = m.Height()
		field       string
	)

	if m.Input.Focused() {
		iHeight := availHeight / 3
		m.Input.SetHeight(iHeight)
		field = m.Input.View()
		availHeight -= iHeight
	}

	m.SetSize(m.width, availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	if m.Input.Focused() {
		sections = append(sections, field)
	}

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
