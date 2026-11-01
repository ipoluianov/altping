package forms

import (
	"image/color"

	"github.com/ipoluianov/nui/ui"
)

// The color themes offered in the settings (config.Settings.Theme)
const (
	themeDark  = ""
	themeLight = "light"
)

// themeColors is a color of the application in the dark and the light theme
type themeColors struct {
	dark, light color.RGBA
}

// get returns the color for the current theme
func (c themeColors) get() color.RGBA {
	if ui.IsDarkTheme {
		return c.dark
	}
	return c.light
}

var (
	colorLink          = themeColors{ui.ColorFromHex("#3d8bf2"), ui.ColorFromHex("#1d6fb5")}
	colorSecondaryText = themeColors{ui.ColorFromHex("#555555"), ui.ColorFromHex("#8c8c8c")}
)

// ApplyTheme switches the application to the theme of the settings and
// repaints the open windows. The colors of the table follow on its next
// update; the colors and icons the application sets itself follow
// ui.OnThemeChanged.
func ApplyTheme(theme string) {
	if theme == themeLight {
		ui.ApplyLightTheme()
	} else {
		ui.ApplyDarkTheme()
	}
}
