package list

import (
	"fmt"

	bubblekey "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type ActionFunc func(items ...*Item) tea.Cmd

type List struct {
	Model         list.Model
	Title         string
	SelectionList bool
	ActionFunc    ActionFunc
	Hash          map[string]string
	Style         list.Styles
	id            int
	style.Frame
	*Items
}

func NewList(title string) *List {
	m := List{
		Frame: style.DefaultFrameStyle(),
		Items: NewItems(),
		Title: title,
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
	l.Styles = style.ListStyles()
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
	m.SetShowKeys()
	m.SetModel()
	return m
}

func (m *List) SetModel() *List {
	m.Items.Process()
	m.Model = NewListModel(m.Width(), m.Height(), m.Items)
	m.Model.Title = m.Title
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

//func (m *List) Update(msg tea.Msg) (*List, tea.Cmd) {
func (m *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case m.SelectionList:
			switch {
			case key.Matches(msg, key.PrevScreen):
				m.SelectionList = false
				cmds = append(cmds, m.ShowVisibleItemsCmd())
			case key.Matches(msg, key.Enter):
				cmds = append(cmds, ReturnSelectionsCmd())
			}
		case m.MultiSelect():
			switch {
			case key.Matches(msg, key.Enter):
				if m.SelectionList {
					cmds = append(cmds, ReturnSelectionsCmd())
				}
				m.SelectionList = true
				cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
			case key.Matches(msg, key.UnToggleAllItems):
				m.Items.DeselectAllItems()
				cmds = append(cmds, m.ShowVisibleItemsCmd())
			case key.Matches(msg, key.ToggleAllItems):
				m.Items.ToggleAllSelectedItems()
				cmds = append(cmds, m.ShowVisibleItemsCmd())
			}
		default:
			switch {
			case key.Matches(msg, key.Enter):
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
		cmds = append(cmds, m.ActionFunc(items...))
	case SortItemsMsg:
		m.Items.SetItems(msg.Items...)
		m.Items.Process()
		cmds = append(cmds, m.ShowVisibleItemsCmd())
	case SetItemsMsg:
		m.Model.SetItems(msg.Items)
	case ToggleItemChildrenMsg:
		switch msg.ShowChildren {
		case true:
			m.Model.Select(msg.Index())
			m.Items.CloseItemList(msg.Index())
		default:
			m.Model.CursorDown()
			m.Items.OpenItemList(msg.Index())
		}
		cmds = append(cmds, m.ShowVisibleItemsCmd())
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	var (
		sections    []string
		availHeight = m.Height()
	)

	m.SetSize(m.Width(), availHeight)
	content := m.Model.View()
	sections = append(sections, content)

	return lipgloss.NewStyle().Height(availHeight).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func ListKeyMap() list.KeyMap {
	km := list.DefaultKeyMap()
	km.NextPage = bubblekey.NewBinding(
		bubblekey.WithKeys("right", "l", "pgdown"),
		bubblekey.WithHelp("l/pgdn", "next page"),
	)
	km.Quit = bubblekey.NewBinding(
		bubblekey.WithKeys("ctrl+c", "esc"),
		bubblekey.WithHelp("ctrl+c", "quit"),
	)
	return km
}

func PrintItems(items ...*Item) tea.Cmd {
	for _, i := range items {
		fmt.Println(i.Content())
	}
	return tea.Quit
}
