package app

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
)

type PageInitializer func(cmpnt.Props) reactea.SomeComponent

type Page struct {
	Data        []teacozy.Items
	Name        string
	CurrentItem int
	Initializer PageInitializer
	*cmpnt.Pager
}

func NewPage(name string, data ...teacozy.Items) *Page {
	page := &Page{
		Data:  data,
		Name:  name,
		Pager: cmpnt.New(),
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

	props := cmpnt.Props{
		Items:      cmpnt.NewItems(p.Data[idx]),
		Current:    p.Current,
		SetCurrent: p.SetCurrent,
	}

	if p.Initializer != nil {
		return p.Initializer(props)
	}

	return p.Pager
}

func (p Page) Current() int {
	return p.CurrentItem
}

func (p *Page) SetCurrent(idx int) {
	p.CurrentItem = idx
}
