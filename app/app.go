package app

import (
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	defaultRoute string

	confirmChoices bool
	readOnly       bool

	width  int
	height int

	numSelected int
	limit       int
	CurrentItem int
	noLimit     bool

	footer string

	choices teacozy.Items
	keyMap  keys.KeyMap

	title  string
	header string

	help keys.KeyMap
}
