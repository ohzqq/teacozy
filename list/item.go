package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/input"
)

type ListType int

const (
	Ul ListType = iota
	Ol
)

// Items holds the values to configure list.DefaultDelegate.
type Items struct {
	list.DefaultDelegate
	ParseFunc       func() []*Item
	ListType        ListType
	width           int
	height          int
	editable        bool
	toggledItems    map[int]struct{}
	limit           int
	prefix          string
	toggledPrefix   string
	untoggledPrefix string
}

// ParseItems is a func to return a slice of Item.
type ParseItems func() []*Item

// Item is a struct to satisfy list.DefaultItem.
type Item struct {
	title       string
	desc        string
	filterValue string
}

// NewItems initializes an Items.
func NewItems(fn ParseItems, opts ...ItemOption) Items {
	items := Items{
		ParseFunc:       fn,
		ListType:        Ol,
		width:           0,
		height:          0,
		editable:        true,
		limit:           0,
		toggledItems:    make(map[int]struct{}),
		prefix:          ">",
		toggledPrefix:   "◉",
		untoggledPrefix: "○",
		DefaultDelegate: list.NewDefaultDelegate(),
	}

	for _, opt := range opts {
		opt(&items)
	}

	return items
}

type ToggleItemMsg struct {
	Index int
}

func ToggleItem(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{Index: idx}
	}
}

// Selectable returns if a list's items can be toggled.
func (m Items) Selectable() bool {
	return m.limit != 0
}

// SetNoLimit allows all items in a list to be toggled.
func (m *Items) SetNoLimit() {
	m.limit = -1
}

// SetSelectNone renders a non-selectable list.
func (m *Items) SetSelectNone() {
	m.limit = 0
}

func (m *Items) SetLimit(n int) {
	m.limit = n
}

// ToggledItems returns the slice of item indices.
func (m Items) ToggledItems() []int {
	var items []int
	for k, _ := range m.toggledItems {
		items = append(items, k)
	}
	return items
}

// MultiSelectable is a convenience method to check if more than one item can be
// toggled.
func (m Items) MultiSelectable() bool {
	if m.limit > 1 {
		return true
	}
	if m.limit == -1 {
		return true
	}
	return false
}

// Limit returns the max number of items that can be toggled.
func (m Items) Limit() int {
	return m.limit
}

func (items *Items) ToggleItems(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ToggleItemMsg:
		idx := msg.Index
		if _, ok := items.toggledItems[idx]; ok {
			delete(items.toggledItems, idx)
		} else {
			no := items.limit
			if items.limit == -1 {
				no = len(m.Items())
			}
			if len(items.toggledItems) < no {
				items.toggledItems[idx] = struct{}{}
			}
		}
		m.CursorDown()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, toggleItem):
			cmd = ToggleItem(m.Index())
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func InsertOrRemoveItems(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case RemoveItemMsg:
		m.RemoveItem(m.Index())

	case InsertItemMsg:
		if msg.Value != "" {
			item := NewItem(msg.Value)
			cmd = m.InsertItem(m.Index()+1, item)
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, input.Reset)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, insertItem):
			cmd = input.Focus
			cmds = append(cmds, cmd)
		case key.Matches(msg, removeItem):
			cmd = RemoveItem(m.Index())
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// ItemOption sets options for Items.
type ItemOption func(*Items)

// ItemsStringSlice returns a ParseItems for a slice of strings.
func ItemsStringSlice(items []string) ParseItems {
	fn := func() []*Item {
		li := make([]*Item, len(items))
		for i, item := range items {
			li[i] = NewItem(item)
		}
		return li
	}
	return fn
}

// ItemsMapSlice returns a ParseItems for a slice of map[string]string.
func ItemsMapSlice(items []map[string]string) ParseItems {
	fn := func() []*Item {
		li := make([]*Item, len(items))
		for i, item := range items {
			for k, v := range item {
				li[i] = NewItem(k).SetDescription(v)
			}
		}
		return li
	}
	return fn
}

// ItemsMap returns a ParseItems for a map[string]string.
func ItemsMap(items map[string]string) ParseItems {
	fn := func() []*Item {
		var li []*Item
		for k, v := range items {
			li = append(li, NewItem(k).SetDescription(v))
		}
		return li
	}
	return fn
}

// NewItem returns an Item struct.
func NewItem(title string) *Item {
	return &Item{
		title:       title,
		desc:        title,
		filterValue: title,
	}
}

// SetDescription sets the Item description.
func (i *Item) SetDescription(desc string) *Item {
	i.desc = desc
	return i
}

// SetFilterValue sets the Item filter value.
func (i *Item) SetFilterValue(val string) *Item {
	i.filterValue = val
	return i
}

func (i Item) FilterValue() string { return i.filterValue }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
