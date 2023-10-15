package app

import (
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
)

type Help struct {
	list  *list.KeyMap
	pager *pager.KeyMap
	items *list.DelegateKeyMap
	app   *KeyMap
}
