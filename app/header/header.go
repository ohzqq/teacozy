package header

import (
	"fmt"
)

type Header string

func RenderHeader(text Header, w, h int) string {
	return fmt.Sprintf("%s", text)
}
