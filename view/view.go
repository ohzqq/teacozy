package view

import "github.com/charmbracelet/bubbles/viewport"

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.Props]

	Viewport viewport.Model
}

type Props struct {
}
