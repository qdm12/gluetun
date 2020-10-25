package constants

import "github.com/fatih/color"

func ColorUnbound() *color.Color {
	return color.New(color.FgCyan)
}

func ColorTinyproxy() *color.Color {
	return color.New(color.FgHiGreen)
}

func ColorOpenvpn() *color.Color {
	return color.New(color.FgHiMagenta)
}
