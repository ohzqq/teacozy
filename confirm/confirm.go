package list

import "github.com/londek/reactea"

type Confirm struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[ConfirmProps]

	question string
}

type ConfirmProps struct {
	Confirm func(bool)
}
