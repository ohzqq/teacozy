package app

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
)

type Route func(*cmpnt.Pager) reactea.SomeComponent

type Page struct {
	Data []teacozy.Items
	Name string
	*cmpnt.Pager
}

func NewPage(name string, data ...teacozy.Items) *Page {
	page := &Page{
		Data:  data,
		Name:  name,
		Pager: cmpnt.New(),
	}
	if len(data) > 0 {
		page.Pager = cmpnt.New(data[0])
	}
	return page
}

func (p *Page) UpdateProps(id string) *cmpnt.Pager {
	idx, err := strconv.Atoi(id)
	if err != nil {
		return p.Pager
	}
	props := p.Pager.NewProps(p.Data[idx])
	p.Pager.Init(props)
	return p.Pager
}

func (p *Page) Initialize(id string, fn Route) reactea.SomeComponent {
	pager := p.UpdateProps(id)
	return fn(pager)
}
