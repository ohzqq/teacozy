package list

import (
	"fmt"
	"io"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/teacozy/input"
	"golang.org/x/exp/slices"
)

// Items holds the values to configure list.DefaultDelegate.
type Items struct {
	list.DefaultDelegate
	ParseFunc func() []*Item
	ListType  list.ListType
	width     int
	height    int
	editable  bool
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
		ParseFunc: fn,
		editable:  true,
		ListType:  list.Ol,
	}
	for _, opt := range opts {
		opt(&items)
	}
	return items
}

// NewDelegate returns a list.ItemDelegate.
func (items Items) NewDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.SetListType(items.ListType)
	return del
}

func (items Items) ItemDelegate() list.ItemDelegate {
	del := items.NewDelegate()

	//if items.editable {
	//  del.UpdateFunc = EditItems
	//  km := list.DefaultKeyMap()
	//  del.ShortHelpFunc = func() []key.Binding {
	//    return []key.Binding{
	//      km.InsertItem,
	//      km.RemoveItem,
	//    }
	//  }
	//}

	items.DefaultDelegate = del
	return items
}

// Render satisfies list.ItemDelegate.
func (d Items) Render(w io.Writer, m list.Model, index int, item list.Item) {

	var (
		prefix     string
		padding    = len(strconv.Itoa(len(m.Items())))
		isSelected = index == m.Index()
	)

	// style prefix
	switch d.ListType {
	case list.Ol:
		p := "%" + strconv.Itoa(padding) + "d."
		prefix = fmt.Sprintf(p, index+1)
	default:
		prefix = " "
	}

	if m.MultiSelectable() {
		if slices.Contains(m.ToggledItems(), index) {
			prefix = fmt.Sprint("[x]" + prefix)
		} else {
			prefix = fmt.Sprint("[ ]" + prefix)
		}
	}

	if isSelected {
		if d.ListType == list.Ul && !m.MultiSelectable() {
			//prefix = m.prefix
			prefix = "> "
		}
		//prefix = s.Prefix.Render(prefix)
		prefix = prefix
	}

	fmt.Fprintf(w, "%s", prefix)
	d.DefaultDelegate.Render(w, m, index, item)
}

func EditItemz(m *Model) func(tea.Msg, *list.Model) tea.Cmd {
	m.hasInput = true
	m.editable = true
	m.SetInput("Insert Item: ", InsertItem)
	//m.ConfigureList(WithFiltering(false), WithLimit(0))
	return EditItems
}

func EditItems(msg tea.Msg, m *list.Model) tea.Cmd {
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
		case key.Matches(msg, m.KeyMap.InsertItem):
			cmd = input.Focus
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.KeyMap.RemoveItem):
			cmd = RemoveItem(m.Index())
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (items Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return EditItems(msg, m)
}

// ItemOption sets options for Items.
type ItemOption func(*Items)

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
