package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	cozykey "github.com/ohzqq/key"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	Model            list.Model
	AllItems         Items
	Items            Items
	Selections       Items
	Keys             cozykey.KeyMap
	Title            string
	ShowSelectedOnly bool
	FocusedView      string
	IsMultiSelect    bool
	width            int
	height           int
	ShowMenu         bool
	frame            lipgloss.Style
}

func New(title string, items Items) *List {
}

func (l List) Width() int {
	return util.TermWidth()
}

func (l List) Height() int {
	return util.TermHeight()
}

//func (l List) GetHeight(items []list.Item) int {
//  max := util.TermHeight()
//  total := len(items)
//  cur := l.Model.Height()

//  switch {
//  case l.isFullScreen:
//    return max
//  case cur > max:
//    return max
//  case total < max:
//    return total + 6
//  default:
//    return max
//  }
//}
