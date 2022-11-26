package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Style struct {
	List   list.Styles
	Item   style.ItemStyle
	Widget style.Frame
	Frame  style.Frame
}

func DefaultStyle() Style {
	tui := Style{
		Item:   style.ItemStyles(),
		List:   style.ListStyles(),
		Widget: DefaultWidgetStyle(),
		Frame:  style.DefaultFrameStyle(),
	}

	return tui
}

func DefaultWidgetStyle() style.Frame {
	return style.Frame{
		MinWidth:  util.TermWidth(),
		MinHeight: util.TermHeight() / 4,
	}
}
