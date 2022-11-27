package tui

import (
	"github.com/ohzqq/teacozy/menu"
)

func NewMenu(k, h string) *menu.Menu {
	m := menu.New(k, h)
	m.Frame = DefaultWidgetStyle()
	return m
}
