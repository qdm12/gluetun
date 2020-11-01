package constants

import "github.com/fatih/color"

func ColorUnbound() *color.Color {
	return color.New(color.FgCyan)
}

func ColorOpenvpn() *color.Color {
	return color.New(color.FgHiMagenta)
}
