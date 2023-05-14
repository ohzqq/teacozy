package app

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
)

type PageInitializer func(cmpnt.Items) reactea.SomeComponent

type Page struct {
	Data []teacozy.Items
	data []cmpnt.Items
	Name string

	CurrentItem int
	selected    map[int]struct{}

	Initializer PageInitializer
	*cmpnt.Pager
}

func NewPage(name string, data ...teacozy.Items) *Page {
	page := &Page{
		Name:     name,
		Pager:    cmpnt.New(),
		selected: make(map[int]struct{}),
	}
	for _, d := range data {
		items := cmpnt.NewItems(d)
		items.Current = page.Current
		items.SetCurrent = page.SetCurrent
		items.IsSelected = page.ItemIsSelected
		page.data = append(page.data, items)
	}
	page.InitFunc(page.Pager.Initializer)
	return page
}

func (p *Page) InitFunc(fn PageInitializer) *Page {
	p.Initializer = fn
	return p
}

func (p *Page) UpdateProps(id string) reactea.SomeComponent {
	idx, err := strconv.Atoi(id)
	if err != nil {
		idx = 0
	}

	if p.Initializer != nil {
		return p.Initializer(p.data[idx])
	}

	return p.Pager
}

func (p Page) SelectedItems() map[int]struct{} {
	return p.selected
}

func (p Page) ItemIsSelected(idx int) bool {
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
