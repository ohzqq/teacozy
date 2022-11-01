package prompt

import "fmt"

var testData = []string{
	"one",
	"two",
	"three",
}

func TestPrompt() *Prompt {
	p := New("test", testData, true)
	c := p.Choose()
	fmt.Printf("choices %v\n", c)
	return p
}
