package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	Model            list.Model
	Title            string
	MultiSelect      bool
	ShowKeys         bool
	ShowSelectedOnly bool
	isForm           bool
	Keys             KeyMap
	Items            Items
	width            int
	height           int
	Style            list.Styles
}

func NewList(title string, items Items) *List {
	w, h := TermSize()
	p := List{
		Items:  items,
		width:  w,
		height: h,
		Keys:   DefaultKeys(),
		Title:  title,
	}
	p.Model = p.InitList()
	return &p
}

func (m *List) InitList() list.Model {
	l := m.Items.List()
	l.SetSize(m.width, m.height)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Title = m.Title

	if m.Title == "" {
		l.Title = ""
	}
	m.Model = l
	return l
}

func (m List) Height() int {
	return TermHeight()
}

func (m *List) SetItems(items Items) *List {
	m.Items = items
	return m
}

func (m *List) SetMultiSelect() *List {
	m.MultiSelect = true
	m.Items.SetMultiSelect()
	del := NewItemDelegate()
	del.MultiSelect()
	m.Model.SetDelegate(del)
	return m
}

func (m *List) SetShowKeys() *List {
	m.ShowKeys = true
	m.Items.SetShowKeys()
	return m
}

func (m *List) SetSize(w, h int) *List {
	m.width = w
	m.height = h
	m.Model.SetSize(w, h)
	return m
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
		switch {
		case key.Matches(msg, m.Keys.Prev):
			m.ShowSelectedOnly = false
			cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
		}
		switch {
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
	case ToggleItemListMsg:
		switch msg.ListOpen {
		case true:
			m.Items.CloseItemList(msg.Index())
		default:
			m.Items.OpenItemList(msg.Index())
		}
		cmds = append(cmds, UpdateVisibleItemsCmd("visible"))
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *List) Init() tea.Cmd {
	m.Model = m.InitList()
	return nil
}

func (m List) View() string {
	return m.Model.View()
}
