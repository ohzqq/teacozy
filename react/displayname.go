package react

import (
	"fmt"
)

// Our prop(s) is a string itself!
type DisplayProps = string

// Stateless components?!?!
func DisplayRenderer(text DisplayProps, width, height int) string {
	return fmt.Sprintf("OMG! Hello %s!", text)
}
