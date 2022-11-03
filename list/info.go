package list

import (
	"github.com/charmbracelet/bubbles/viewport"
)

type Info struct {
	model   viewport.Model
	Content string
}

func NewInfo(w, h int) Info {
	info := Info{
		model: viewport.New(w, h),
	}
	return info
}

func (i *Info) SetContent(content string) {
	i.model.SetContent(content)
}

func (i *Info) SetSize(w, h int) {
	i.model = viewport.New(w, h)
}
