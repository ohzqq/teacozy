package list

import (
	"github.com/charmbracelet/bubbles/list"
)

type Selections struct {
	items       Items
	indices     []int
	ReturnItems bool
}

func (s Selections) HasItems() bool {
	return len(s.items) > 0
}

func (s Selections) Items() Items {
	var items Items
	for _, i := range s.items {
		item := i.(Item)
		item.isSelected = true
		items = append(items, item)
	}
	return items
}

func (l *Model) CurrentlySelectedItems() ([]list.Item, []int) {
	var items []list.Item
	var indices []int
	for idx, i := range l.Model.Items() {
		switch li := i.(type) {
		case Item:
			if li.IsSelected() {
				items = append(items, i)
				indices = append(indices, idx)
			}
		}
	}
	return items, indices
}

func (l *Model) SelectAll() {
	for idx, it := range l.Model.Items() {
		i := it.(Item)
		i.Toggle()
		l.Model.SetItem(idx, i)
	}
}
