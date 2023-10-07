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
	ListType        list.ListType
	prefix          string
	toggledPrefix   string
	untoggledPrefix string
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
	case list.Ol:
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
		if d.ListType == list.Ul && !d.MultiSelectable() {
			prefix = d.prefix
		}
		//prefix = s.Prefix.Render(prefix)
		prefix = prefix
	}

	fmt.Fprintf(w, "%s", prefix)
	d.DefaultDelegate.Render(w, m, index, item)
}
