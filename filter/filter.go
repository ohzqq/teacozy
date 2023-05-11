package filter

import (
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/state"
)

type Page struct {
	Pager cmpnt.PagerProps
	Input cmpnt.InputProps
}

func New(items teacozy.Items) *teacozy.Page {

	input := cmpnt.NewTextInput()
	input.Init(cmpnt.InputProps{SetValue: state.SetInputValue})

	pager := cmpnt.NewPager()
	props := cmpnt.NewPagerProps(items)
	props.InputValue = input.Input.Value
	pager.Init(props)
	page := teacozy.NewPage("filter", pager)
	page.SetHeader(input)
	return page
}
