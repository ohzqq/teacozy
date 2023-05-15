package cmpnt

import "github.com/ohzqq/teacozy"

func NewPage(name string, data ...teacozy.Items) *teacozy.Page {
	page := teacozy.NewPage(name, data...)
	p := New()
	p.Init(page)
	page.InitFunc(p.Initializer).Update()
	return page
}
