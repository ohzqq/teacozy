package form

import "github.com/ohzqq/teacozy/list"

type Model struct {
	list.List
	Fields       Fields
	CurrentField *Field
}

func New(title string) Model {
	return Model{
		Fields: make(map[string]*Field),
	}
}
