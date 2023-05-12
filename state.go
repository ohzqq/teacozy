package teacozy

import (
	"github.com/ohzqq/teacozy/util"
)

type Page struct {
	numSelected int
	limit       int
	noLimit     bool
	current     int
	selected    map[int]struct{}

	width  int
	height int

	confirmChoices bool
	readOnly       bool

	inputValue string

	style PageStyles
}

type PageStyles struct {
	List Style
}

func NewPage() *Page {
	c := &Page{
		limit:    10,
		selected: make(map[int]struct{}),
		width:    util.TermWidth(),
		height:   util.TermHeight() - 2,
		style: PageStyles{
			List: DefaultStyle(),
		},
	}

	return c
}

func (p Page) CurrentItem() int {
	return p.current
}

func (p Page) Width() int {
	return p.width
}

func (p Page) Height() int {
	return p.height
}

func (p Page) Limit() int {
	return p.limit
}

func (p Page) NoLimit() bool {
	return p.noLimit
}

func (p Page) ReadOnly() bool {
	return p.readOnly
}

func (p Page) Style() PageStyles {
	return p.style
}

func (p *Page) SetLimit(n int) *Page {
	p.limit = n
	return p
}

func (p *Page) SetWidth(n int) *Page {
	p.width = n
	return p
}

func (p *Page) SetHeight(n int) *Page {
	p.height = n
	return p
}

func (p *Page) SetCurrent(n int) *Page {
	p.current = n
	return p
}

func (p *Page) SetNoLimit(n bool) *Page {
	p.noLimit = n
	return p
}

func (p *Page) SetReadOnly(n bool) *Page {
	p.readOnly = n
	return p
}

func (p *Page) SetInputValue(val string) *Page {
	p.inputValue = val
	return p
}
