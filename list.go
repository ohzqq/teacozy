package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List struct {
	Model            list.Model
	Input            textarea.Model
	Title            string
	MultiSelect      bool
	ShowKeys         bool
	ShowSelectedOnly bool
	Editable         bool
	FormChanged      bool
	Keys             KeyMap
	SaveFormFunc     SaveFormFunc
	Hash             map[string]string
	Style            list.Styles
	id               int
	Frame
	Items
}

func NewList(title string, items Items) *List {
	m := DefaultList()
	m.SetItems(items)
	m.SetTitle(title)
	m.Model = ListModel(m.Width(), m.Height(), m.Items)
	return m.ChooseOne()
}

func DefaultList() *List {
	m := List{
		Keys:         DefaultKeys(),
		SaveFormFunc: SaveFormAsHashCmd,
		Frame:        DefaultFrameStyle(),
		Items:        NewItems(),
	}
	m.Frame.MinHeight = 10
	return &m
}

func ListModel(w, h int, items Items) list.Model {
	l := list.New(items.Visible(), items.Delegate, w, h)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Styles = ListStyles()
	return l
}

func (m *List) ChooseOne() *List {
	m.Model = m.Items.List()
	return m
}

func (m *List) ChooseMany() *List {
	m.Model = m.Items.List()
	m.SetMultiSelect()
	return m
}

func (m *List) Edit() *List {
	m.Model = m.Items.List()
	m.SetShowKeys()
	m.Editable = true
	return m
}

func (m *List) SetItems(items Items) *List {
	m.Items = items
	m.Items.Process()
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
	m.Items.Delegate.MultiSelect()
	m.Model.SetDelegate(m.Items.Delegate)
	m.MultiSelect = true
	return m
}

func (m *List) SetShowKeys() *List {
	m.Items.Delegate.ShowKeys()
	m.Model.SetDelegate(m.Items.Delegate)
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
			if key.Matches(msg, Keys.SaveAndExit) {
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
			case key.Matches(msg, m.Keys.Prev):
				m.ShowSelectedOnly = false
				cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
			}
			switch {
			case m.Editable:
				switch {
				case key.Matches(msg, Keys.SaveAndExit):
					cmds = append(cmds, FormChangedCmd())
				}
			case m.MultiSelect:
				switch {
				case key.Matches(msg, Enter):
					if m.ShowSelectedOnly {
						cmds = append(cmds, ReturnSelectionsCmd())
					}
					m.ShowSelectedOnly = true
					cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
				case key.Matches(msg, m.Keys.SelectAll):
					m.Items.ToggleAllSelectedItems()
					cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
				}
			default:
				switch {
				case key.Matches(msg, m.Keys.Enter):
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
		items := m.Items.Display(string(msg))
		m.Model.SetItems(items)
	case ToggleSelectedItemMsg:
		m.Items.ToggleSelectedItem(msg.Index())
	case ReturnSelectionsMsg:
		cmds = append(cmds, tea.Quit)
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
		cmds = append(cmds, SetItemsCmd(msg.Items))
	case SetItemsMsg:
		l := msg.Items.List()
		l.SetSize(m.Width(), m.Height())
		l.Title = m.Title
		m.Model = l
	case ToggleItemListMsg:
		switch msg.ListOpen {
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
