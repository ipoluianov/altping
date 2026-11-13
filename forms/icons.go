package forms

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/ipoluianov/nui/ui"
)

// icons/*.png are rendered at 32x32 from the matching icons/*.svg
//
//go:embed icons/*.png
var iconsFS embed.FS

// The icons are drawn in one light color for the dark theme; the classic
// light theme shows them in this dark one, the other styles in their own
// (Style.Icon)
var iconColorLight = ui.ColorFromHex("#455a64")

// loadIcon returns the icon in the colors of the current style
func loadIcon(name string) image.Image {
	data, err := iconsFS.ReadFile("icons/" + name + ".png")
	if err != nil {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	if currentStyle.Icon.A == 0 {
		return img
	}
	return tintIcon(img, currentStyle.Icon)
}

// setIcon sets the icon with set now and again on every theme change
func setIcon(name string, set func(img image.Image)) {
	set(loadIcon(name))
	ui.OnThemeChanged.Connect(func(bool) { set(loadIcon(name)) })
}

// tintIcon paints a one-color icon in col, keeping its transparency
func tintIcon(img image.Image, col color.RGBA) image.Image {
	src := image.NewNRGBA(img.Bounds())
	draw.Draw(src, src.Rect, img, img.Bounds().Min, draw.Src)
	for i := 0; i < len(src.Pix); i += 4 {
		src.Pix[i], src.Pix[i+1], src.Pix[i+2] = col.R, col.G, col.B
	}
	return src
}
