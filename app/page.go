package app

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
)

type PageInitializer func(*cmpnt.Pager, teacozy.Items) reactea.SomeComponent

type Page struct {
	Data        []teacozy.Items
	Name        string
	Initializer PageInitializer
	*cmpnt.Pager
}

func NewPage(name string, data ...teacozy.Items) *Page {
	page := &Page{
		Data:  data,
		Name:  name,
		Pager: cmpnt.New(),
	}
	return page
}

func (p *Page) InitFunc(fn PageInitializer) *Page {
	p.Initializer = fn
	return p
}

func (p *Page) UpdateProps(id string) reactea.SomeComponent {
	idx, err := strconv.Atoi(id)
	if err != nil {
		return p.Pager
	}

	if p.Initializer != nil {
		return p.Initializer(p.Pager, p.Data[idx])
	}

	return p.Pager
}
