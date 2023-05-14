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
		Data:     data,
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

	//props := cmpnt.NewItems(p.Data[idx])
	//props.Current = p.Current
	//props.SetCurrent = p.SetCurrent
	//props.Selected = p.SelectedItems

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
