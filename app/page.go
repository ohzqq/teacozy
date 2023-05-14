package app

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/keys"
)

type PageInitializer func(cmpnt.PageProps) cmpnt.Page

type Page struct {
	Data        []teacozy.Items
	Name        string
	CurrentPage int

	CurrentItem int
	selected    map[int]struct{}

	Initializer PageInitializer
	keymap      keys.KeyMap
}

func NewPage(name string, data ...teacozy.Items) *Page {
	page := &Page{
		Data:     data,
		Name:     name,
		keymap:   keys.DefaultKeyMap(),
		selected: make(map[int]struct{}),
	}
	page.InitFunc(cmpnt.NewPage())
	return page
}

func (p *Page) InitFunc(fn PageInitializer) *Page {
	p.Initializer = fn
	return p
}

func (p Page) KeyMap() keys.KeyMap {
	return p.keymap
}

func (p *Page) UpdateProps(id string) reactea.SomeComponent {
	idx, err := strconv.Atoi(id)
	if err != nil {
		idx = 0
	}
	p.CurrentPage = idx

	page := p.Initializer(p)
	p.keymap = page.KeyMap()

	return page.Mount()
}

func (p Page) Items() teacozy.Items {
	return p.Data[p.CurrentPage]
}

func (p Page) SelectedItems() map[int]struct{} {
	return p.selected
}

func (p Page) IsSelected(idx int) bool {
	if _, ok := p.selected[idx]; ok {
		return true
	}
	return false
}

func (p Page) Current() int {
	return p.CurrentItem
}

func (p *Page) SetCurrent(idx int) {
	p.CurrentItem = idx
}

func (c *Page) SelectItem(idx int) {
	c.selected[idx] = struct{}{}
}

func (c *Page) DeselectItem(idx int) {
	delete(c.selected, idx)
}
