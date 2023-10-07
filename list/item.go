package list

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/input"
	"golang.org/x/exp/slices"
)

type ListType int

const (
	Ul ListType = iota
	Ol
)

// String returns a human-readable string of the list type.
func (f ListType) String() string {
	return [...]string{
		"unordered list",
		"ordered list",
	}[f]
}

// Items holds the values to configure list.DefaultDelegate.
type Items struct {
	list.DefaultDelegate
	li              []list.Item
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
	styles          DefaultItemStyles
}

// ParseItems is a func to return a slice of Item.
type ParseItems func() []*Item

// Item is a struct to satisfy list.DefaultItem.
type Item struct {
	title       string
	desc        string
	filterValue string
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

// Satisfy list.DefaultItem interface
func (i Item) FilterValue() string { return i.filterValue }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }

// NewItems initializes an Items.
func NewItems(fn ParseItems) *Items {
	var li []list.Item
	for _, i := range fn() {
		di := list.DefaultItem(i)
		li = append(li, di)
	}
	items := Items{
		li:              li,
		ParseFunc:       fn,
		ListType:        Ul,
		width:           10,
		height:          10,
		editable:        true,
		limit:           0,
		toggledItems:    make(map[int]struct{}),
		prefix:          ">",
		toggledPrefix:   "x",
		untoggledPrefix: " ",
		styles:          NewDefaultItemStyles(),
	}
	del := list.NewDefaultDelegate()
	del.Styles = items.styles.DefaultItemStyles
	del.ShowDescription = false
	del.SetHeight(1)
	items.DefaultDelegate = del
	//items.DefaultDelegate = list.NewDefaultDelegate()
	//items.Styles = items.styles.DefaultItemStyles
	return &items
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

func (m Items) Chosen() []*Item {
	var ch []*Item
	for _, c := range m.ToggledItems() {
		ch = append(ch, m.li[c].(*Item))
	}
	return ch
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

func (items *Items) ToggleItem(idx int) tea.Cmd {
	return func() tea.Msg {
		if _, ok := items.toggledItems[idx]; ok {
			delete(items.toggledItems, idx)
		} else {
			no := items.limit
			if items.limit == -1 {
				no = len(items.li)
			}
			if len(items.toggledItems) < no {
				items.toggledItems[idx] = struct{}{}
			}
		}
		return ToggleItemMsg{}
	}
}

func (items *Items) UpdateItemToggle(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ToggleItemMsg:
		m.CursorDown()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, toggleItem):
			cmd = items.ToggleItem(m.Index())
			cmds = append(cmds, cmd)
		}
		switch msg.Type {
		case tea.KeyEnter:
			if !items.MultiSelectable() {
				cmd = items.ToggleItem(m.Index())
				cmds = append(cmds, cmd)
			}
			cmds = append(cmds, ChooseItems)
		}
	}
	return tea.Batch(cmds...)
}

func (items *Items) InsertOrRemoveItems(msg tea.Msg, m *list.Model) tea.Cmd {
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

func (items *Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if items.editable {
		cmd = items.InsertOrRemoveItems(msg, m)
		cmds = append(cmds, cmd)
	}

	if items.Selectable() {
		cmd = items.UpdateItemToggle(msg, m)
		cmds = append(cmds, cmd)
	}

	items.li = m.Items()

	return tea.Batch(cmds...)
}

// func (items Items) Height() int  { return 1 }
func (items Items) Spacing() int { return 0 }

// Render satisfies list.ItemDelegate.
func (d Items) Render(w io.Writer, m list.Model, index int, item list.Item) {

	var (
		prefix     = " "
		padding    = len(strconv.Itoa(len(m.Items())))
		isSelected = index == m.Index()
		isToggled  = slices.Contains(d.ToggledItems(), index)
	)

	//style prefix
	if d.MultiSelectable() {
		if isToggled {
			prefix = d.toggledPrefix
		}
	}

	if d.ListType == Ol {
		p := fmt.Sprint("%", padding, "d")
		prefix = fmt.Sprintf(p, index+1)
	}

	if isToggled {
		prefix = d.styles.Toggled.Render(prefix)
	}

	if isSelected {
		switch d.MultiSelectable() {
		case true:
			if d.ListType == Ul {
				prefix = d.toggledPrefix
			}
		default:
			if d.ListType == Ul {
				prefix = d.prefix
			}
		}
		prefix = d.styles.Prefix.Render(prefix)
		//prefix = prefix
	}

	fmt.Fprintf(w, "[%s]", prefix)
	// fmt.Fprintf(w, "%s", item.(list.DefaultItem).Title())
	d.DefaultDelegate.Render(w, m, index, item)
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

type ToggleItemMsg struct {
	Index int
}

func ToggleItem(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{Index: idx}
	}
}

type ItemsChosenMsg struct{}

func ChooseItems() tea.Msg {
	return ItemsChosenMsg{}
}

// InsertItemMsg holds the title of the item to be inserted.
type InsertItemMsg struct {
	Value string
}

// InsertItem returns a tea.Cmd to insert an item into a list.
func InsertItem(val string) tea.Cmd {
	return func() tea.Msg {
		return InsertItemMsg{
			Value: val,
		}
	}
}

// RemoveItemMsg is a struct for the index to be removed.
type RemoveItemMsg struct {
	Index int
}

// RemoveItem returns a tea.Cmd for removing the item at index n.
func RemoveItem(idx int) tea.Cmd {
	return func() tea.Msg {
		return RemoveItemMsg{Index: idx}
	}
}
