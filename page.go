package teacozy

import (
	"strconv"

	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type PageInitializer func(PageProps) PageComponent

type Page struct {
	reactea.BasicComponent
	Data        []Items
	Name        string
	CurrentPage int

	CurrentItem int
	selected    map[int]struct{}

	Initializer PageInitializer
	keymap      keys.KeyMap
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

func NewPageComponent(title string, items ...Items) func(PageProps) PageComponent {
	return func(props PageProps) PageComponent {
		return NewPage(title, items...)
	}
}

func NewPage(name string, data ...Items) *Page {
	page := &Page{
		Data:     data,
		Name:     name,
		keymap:   keys.DefaultKeyMap(),
		selected: make(map[int]struct{}),
	}
	//page.InitFunc(NewPageComponent(name, data...))
	return page
}

func (p *Page) InitFunc(fn PageInitializer) *Page {
	p.Initializer = fn
	return p
}

func (p Page) KeyMap() keys.KeyMap {
	return p.keymap
}

func (p *Page) Mount() reactea.SomeComponent {
	return p
}

func (p Page) Render(w, h int) string {
	return p.Mount().Render(w, h)
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

func (p Page) Items() Items {
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
