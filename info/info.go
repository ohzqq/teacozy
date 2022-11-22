package info

import "github.com/charmbracelet/bubbles/viewport"

type Info struct {
	Model     viewport.Model
	HideKeys  bool
	IsVisible bool
	Editable  bool
	content   []string
	//Style     FieldStyle
	//Frame     Frame
	Data   FormData
	Fields *Fields
}
