package list

import (
	"fmt"
	"io"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/bubbles/list"
	"golang.org/x/exp/slices"
)

type ItemDelegate struct {
	list.DefaultDelegate
	ListType        ListType
	prefix          string
	toggledPrefix   string
	untoggledPrefix string
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

func (items Items) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if items.editable {
		cmd = InsertOrRemoveItems(msg, m)
		cmds = append(cmds, cmd)
	}

	if items.Selectable() {
		cmd = items.ToggleItems(msg, m)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
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
	case Ol:
		p := "%" + strconv.Itoa(padding) + "d."
		prefix = fmt.Sprintf(p, index+1)
	default:
		prefix = " "
	}

	if d.MultiSelectable() {
		if slices.Contains(d.ToggledItems(), index) {
			prefix = fmt.Sprint("[x]" + prefix)
		} else {
			prefix = fmt.Sprint("[ ]" + prefix)
		}
	}

	if isSelected {
		if d.ListType == Ul && !d.MultiSelectable() {
			prefix = d.prefix
		}
		//prefix = s.Prefix.Render(prefix)
		prefix = prefix
	}

	fmt.Fprintf(w, "%s", prefix)
	d.DefaultDelegate.Render(w, m, index, item)
}
