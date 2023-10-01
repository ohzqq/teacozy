package pager

import (
	"log"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	text := strings.Join(choiceSlice, "\n- ")
	m := New(text, RenderMarkdown)
	m.SetSize(0, 10)

	p := tea.NewProgram(m)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

}

var choiceSlice = []string{
	"Artichoke",
	"Baking Flour",
	"Bananas",
	"Barley",
	"Bean Sprouts",
	"Bitter Melon",
	"Blood Orange",
	"Brown Sugar",
	"Cashew Apple",
	"Cashews",
	"Cat Food",
	"Coconut Milk",
	"Cucumber",
	"Curry Paste",
	"Currywurst",
	"Dill",
	"Dragonfruit",
	"Dried Shrimp",
	"Eggs",
	"Fish Cake",
	"Furikake",
	"Garlic",
	"Gherkin",
	"Ginger",
	"Granulated Sugar",
	"Grapefruit",
	"Green Onion",
	"Hazelnuts",
	"Heavy whipping cream",
	"Honey Dew",
	"Horseradish",
	"Jicama",
	"Kohlrabi",
	"Leeks",
	"Lentils",
	"Licorice Root",
	"Meyer Lemons",
	"Milk",
	"Molasses",
	"Muesli",
	"Nectarine",
	"Niagamo Root",
	"Nopal",
	"Nutella",
	"Oat Milk",
	"Oatmeal",
	"Olives",
	"Papaya",
	"Party Gherkin",
	"Peppers",
	"Persian Lemons",
	"Pickle",
	"Pineapple",
	"Plantains",
	"Pocky",
	"Powdered Sugar",
	"Quince",
	"Radish",
	"Ramps",
	"Star Anise",
	"Sweet Potato",
	"Tamarind",
	"Unsalted Butter",
	"Watermelon",
	"Weißwurst",
	"Yams",
	"Yeast",
	"Yuzu",
}
