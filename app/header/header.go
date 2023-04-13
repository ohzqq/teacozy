package header

import (
	"fmt"

	"github.com/londek/reactea"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
}

type Props string

func (c Component) Render(int, int) string {
	return fmt.Sprintf("%s", c.Props())
}
