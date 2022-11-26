package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/form"
	"github.com/ohzqq/teacozy/style"
)

type Item struct {
	idx          int
	IsHidden     bool
	IsSelected   bool
	MultiSelect  bool
	ShowChildren bool
	ShowKey      bool
	ListLevels   int
	Level        int
	Parent       *Item
	Children     *Items
	hasFields    bool
	style        style.ItemStyle
	Meta         *form.Fields
	teacozy.Field
}

func NewItem(item teacozy.Field) *Item {
	return &Item{
		Field:    item,
		Meta:     form.NewFields(),
		Children: NewItems(),
		style:    style.ItemStyles(),
	}
}

func (i *Item) SetMeta(meta teacozy.Fields) {
	i.Meta.SetData(meta)
}

func (i Item) HasMeta() bool {
	return len(i.Meta.Keys()) > 0
}

// Satisfy Fields interface
func (i Item) Get(key string) teacozy.Field {
	return i.Meta.Get(key)
}

func (i Item) Keys() []string {
	return i.Meta.Keys()
}

// Satisfy list.Item interface
func (i Item) FilterValue() string {
	return i.Content()
}

// Item methods
func (i Item) Prefix() string {
	prefix := dash
	if i.MultiSelect {
		prefix = uncheck
		if i.IsSelected {
			prefix = check
		}
	}

	if i.HasChildren() {
		prefix = openSub
		if i.ShowChildren {
			prefix = closeSub
		}
	}

	return prefix
}

func (i Item) Style(current bool) lipgloss.Style {
	var lip lipgloss.Style
	switch {
	case i.IsSelected:
		lip = i.style.Selected
	case i.IsSub():
		lip = i.style.Sub
	default:
		lip = i.style.Normal
	}
	if current {
		lip = i.style.Current
	}
	return lip.Copy().Margin(0, 1, 0, i.Level)
}

func (i Item) Render(current, key bool) string {
	var (
		s    = i.Style(current).Copy().Margin(0, 1, 0, i.Level)
		p    = i.Prefix()
		text = i.Content()
	)

	return s.Render(p + text)
}

func (i Item) Index() int {
	return i.idx
}

func (i *Item) SetMultiSelect() *Item {
	i.MultiSelect = true
	return i
}

func (i *Item) ToggleSelected() *Item {
	i.IsSelected = !i.IsSelected
	return i
}

func (i *Item) Select() *Item {
	i.IsSelected = true
	return i
}

func (i *Item) Deselect() *Item {
	i.IsSelected = false
	return i
}

func (i *Item) Show() *Item {
	i.IsHidden = false
	return i
}

func (i *Item) Hide() *Item {
	i.IsHidden = true
	return i
}

// Sublist methods
func (i *Item) Flatten() []*Item {
	var items []*Item
	depth := 0
	if i.HasChildren() {
		depth++
		for _, item := range i.Children.flat {
			if i.MultiSelect {
				item.SetMultiSelect()
			}
			item.IsHidden = true
			items = append(items, item)
			if item.HasChildren() {
				depth++
				items = append(items, item.Flatten()...)
			}
		}
	}
	i.ListLevels = depth
	return items
}

func (i Item) ListDepth() int {
	depth := 0
	if i.HasChildren() {
		depth++
		for _, item := range i.Children.flat {
			if item.HasChildren() {
				depth++
			}
		}
	}
	return depth
}

func (i Item) HasChildren() bool {
	has := len(i.Children.flat) > 0
	return has
}

func (i Item) TotalChildren() int {
	if i.HasChildren() {
		return len(i.Children.flat)
	}
	return 0
}

func (i *Item) ToggleList() *Item {
	i.ShowChildren = !i.ShowChildren
	return i
}

func (i *Item) Open() *Item {
	i.ShowChildren = true
	return i
}

func (i *Item) Close() *Item {
	i.ShowChildren = false
	return i
}

func (i *Item) SetLevel(l int) *Item {
	i.Level = l
	return i
}

func (i *Item) IsSub() bool {
	return i.Level != 0
}
