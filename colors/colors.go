package colors

import (
	"github.com/fatih/color"
)

var (
	GREEN  *color.Color
	RED    *color.Color
	BLUE   *color.Color
	YELLOW *color.Color
	CYAN   *color.Color
)

// InitColors function create all colors for output
func InitColors() {
	GREEN = color.New(color.FgGreen)
	RED = color.New(color.FgRed)
	BLUE = color.New(color.FgBlue)
	YELLOW = color.New(color.FgYellow)
	CYAN = color.New(color.FgCyan)
}
