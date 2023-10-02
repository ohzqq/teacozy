package list

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/input"
)

// Items holds the values to configure list.DefaultDelegate.
type Items struct {
	list.DefaultDelegate
	ParseFunc func() []*Item
	ListType  list.ListType
	width     int
	height    int
}

// ParseItems is a func to return a slice of Item.
type ParseItems func() []*Item

// Item is a struct to satisfy list.DefaultItem.
type Item struct {
	title       string
	desc        string
	filterValue string
}

// Render satisfies list.ItemDelegate.
func (items Items) Render(w io.Writer, m list.Model, index int, item list.Item) {
	items.DefaultDelegate.Render(w, m, index, item)
}

func (items Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case InsertItemMsg:
		if msg.Value != "" {
			item := NewItem(msg.Value)
			cmd = m.InsertItem(m.Index()+1, item)
			cmds = append(cmds, cmd)
		}
		//cmds = append(cmds, m.input.Reset)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.InsertItem):
			cmds = append(cmds, input.Focus)
			//if m.hasInput {
			//  m.SetShowInput(true)
			//  cmds = append(cmds, m.input.Focus())
			//}
		case key.Matches(msg, m.KeyMap.RemoveItem):
			m.RemoveItem(m.Index())
		}
	}
	return tea.Batch(cmds...)
}

// ItemOption sets options for Items.
type ItemOption func(*Items)

// NewItems initializes an Items.
func NewItems(fn ParseItems, opts ...ItemOption) Items {
	items := Items{
		ParseFunc: fn,
	}
	for _, opt := range opts {
		opt(&items)
	}
	return items
}

// NewDelegate returns a list.DefaultDelegate.
// func (items Items) NewDelegate() list.DefaultDelegate {
func (items Items) NewDelegate() list.ItemDelegate {
	del := list.NewDefaultDelegate()
	del.SetListType(items.ListType)
	items.DefaultDelegate = del
	return items
}

// OrderedList sets the list.DefaultDelegate ListType.
func OrderedList() ItemOption {
	return func(items *Items) {
		items.ListType = list.Ol
	}
}

// UnrderedList sets the list.DefaultDelegate ListType.
func UnorderedList() ItemOption {
	return func(items *Items) {
		items.ListType = list.Ul
	}
}

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
