package constants

import "github.com/fatih/color"

func ColorUnbound() color.Color {
	return *color.New(color.FgCyan)
}

func ColorTinyproxy() color.Color {
	return *color.New(color.FgHiMagenta)
}

func ColorShadowsocks() color.Color {
	return *color.New(color.FgHiYellow)
}

func ColorShadowsocksError() color.Color {
	return *color.New(color.FgHiRed)
}

func ColorOpenvpn() color.Color {
	return *color.New(color.FgHiMagenta)
}
