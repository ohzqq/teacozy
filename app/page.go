package app

import (
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
)

type Page struct {
	Data []teacozy.Items
	list *cmpnt.Pager
}

func NewPage(data ...teacozy.Items) *Page {
	page := &Page{
		Data: data,
	}
	if len(data) > 0 {
		page.list = cmpnt.New(data[0])
	}
	return page
}
