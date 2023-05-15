package teacozy

import (
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type PageInitializer func(PageProps) PageComponent

type Page struct {
	Data        []Items
	Endpoint    string
	CurrentPage int

	CurrentItem int
	selected    map[int]struct{}
	SetCurPage  func() int

	Initializer PageInitializer
	keyMap      keys.KeyMap
}

type PageProps interface {
	SetCurrent(int)
	Current() int
	IsSelected(int) bool
	Items() Items
}

type PageComponent interface {
	Mount() reactea.SomeComponent
	KeyMap() keys.KeyMap
}

func NewPage(name string, data ...Items) *Page {
	page := &Page{
		Data:     data,
		Endpoint: name,
		keyMap:   keys.DefaultKeyMap(),
		selected: make(map[int]struct{}),
	}
	return page
}

func (p *Page) InitFunc(fn PageInitializer) *Page {
	p.Initializer = fn
	return p
}

func (p Page) KeyMap() keys.KeyMap {
	return p.keyMap
}

func (p *Page) AddItems(data ...Items) {
	p.Data = append(p.Data, data...)
}

func (p *Page) Update() reactea.SomeComponent {
	page := p.Initializer(p)
	p.keyMap = page.KeyMap()
	return page.Mount()
}

func (p Page) Items() Items {
	return p.Data[p.CurrentPage]
}

func (p *Page) SetCurrentPage(idx int) {
	p.CurrentPage = idx
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
