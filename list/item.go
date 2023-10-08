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

const (
	SelectAll = iota - 1
	SelectNone
	SelectOne
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
	KeyMap          DelegateKeyMap
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
		li = append(li, list.DefaultItem(i))
	}
	items := Items{
		li:              li,
		ListType:        Ul,
		toggledItems:    make(map[int]struct{}),
		prefix:          ">",
		toggledPrefix:   "x",
		untoggledPrefix: " ",
		styles:          NewDefaultItemStyles(),
		KeyMap:          DefaultDelegateKeyMap(),
	}
	items.KeyMap.ToggleItem.SetEnabled(false)
	items.SetEditable(items.Editable())
	return &items
}

// NewDelegate returns a list.DefaultDelegate with the default style.
func (items *Items) NewDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.Styles = items.styles.DefaultItemStyles
	return del
}

// Editable returns if a list is editable.
func (items Items) Editable() bool {
	return items.editable
}

func (items *Items) SetEditable(edit bool) {
	items.editable = edit
	items.KeyMap.InsertItem.SetEnabled(edit)
	items.KeyMap.RemoveItem.SetEnabled(edit)
}

// Selectable returns if a list's items can be toggled.
func (items Items) Selectable() bool {
	return items.limit != SelectNone
}

// SetLimit sets the number of choices for a selectable list.
func (items *Items) SetLimit(n int) {
	if n != SelectNone {
		items.KeyMap.ToggleItem.SetEnabled(true)
	}
	items.limit = n
}

// Chosen returns the toggled items.
func (items Items) Chosen() []*Item {
	var ch []*Item
	for _, c := range items.ToggledItems() {
		ch = append(ch, items.li[c].(*Item))
	}
	return ch
}

func (items Items) Len() int {
	return len(items.li)
}

// ToggledItems returns the slice of item indices.
func (i Items) ToggledItems() []int {
	var items []int
	for k, _ := range i.toggledItems {
		items = append(items, k)
	}
	return items
}

// MultiSelectable is a convenience method to check if more than one item can be
// toggled.
func (items Items) MultiSelectable() bool {
	if items.limit > SelectOne {
		return true
	}
	if items.limit == SelectAll {
		return true
	}
	return false
}

// Limit returns the max number of items that can be toggled.
func (items Items) Limit() int {
	return items.limit
}

// ToggleItem toggles the item at index and returns a tea.Cmd.
func (items *Items) ToggleItem(idx int) tea.Cmd {
	return func() tea.Msg {
		if _, ok := items.toggledItems[idx]; ok {
			delete(items.toggledItems, idx)
		} else {
			no := items.limit
			if items.limit == SelectAll {
				no = len(items.li)
			}
			if len(items.toggledItems) < no {
				items.toggledItems[idx] = struct{}{}
			}
		}
		return ToggleItemMsg{}
	}
}

// UpdateItemToggle provides the updates for a selectable list.
func (items *Items) UpdateItemToggle(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ToggleItemMsg:
		m.CursorDown()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, items.KeyMap.ToggleItem):
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

// InsertOrRemoveItems provides the updates for an editable list.
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
		case key.Matches(msg, items.KeyMap.InsertItem):
			cmd = input.Focus
			cmds = append(cmds, cmd)
		case key.Matches(msg, items.KeyMap.RemoveItem):
			cmd = RemoveItem(m.Index())
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update is the update loop for list.ItemDelegate.
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

// Spacing sets the list.DefaultDelegate default spacing to 0.
func (items Items) Spacing() int { return 0 }

// Render satisfies list.ItemDelegate.
func (items Items) Render(w io.Writer, m list.Model, index int, item list.Item) {

	var (
		prefix     = " "
		padding    = len(strconv.Itoa(len(m.Items())))
		isSelected = index == m.Index()
		isToggled  = slices.Contains(items.ToggledItems(), index)
	)

	//style prefix
	if items.MultiSelectable() {
		if isToggled {
			prefix = items.toggledPrefix
		}
	}

	if items.ListType == Ol {
		p := fmt.Sprint("%", padding, "d")
		prefix = fmt.Sprintf(p, index+1)
	}

	if isToggled {
		prefix = items.styles.Toggled.Render(prefix)
	}

	if isSelected {
		switch items.MultiSelectable() {
		case true:
			if items.ListType == Ul {
				prefix = items.toggledPrefix
			}
		default:
			if items.ListType == Ul {
				prefix = items.prefix
			}
		}
		prefix = items.styles.Prefix.Render(prefix)
	}

	fmt.Fprintf(w, "[%s]", prefix)
	items.DefaultDelegate.Render(w, m, index, item)
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
