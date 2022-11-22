package form

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/ohzqq/teacozy/list"
)

type Form struct {
	*Fields
	Model list.Model
	Input textarea.Model
}
