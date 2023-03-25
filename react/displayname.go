package react

import (
	"fmt"
)

// Our prop(s) is a string itself!
type Props = string

// Stateless components?!?!
func Renderer(text Props, width, height int) string {
	return fmt.Sprintf("OMG! Hello %s!", text)
}
