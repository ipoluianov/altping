package forms

import (
	"bytes"
	"embed"
	"image"
	"image/png"
)

// icons/*.png are rendered at 32x32 from the matching icons/*.svg
//
//go:embed icons/*.png
var iconsFS embed.FS

func loadIcon(name string) image.Image {
	data, err := iconsFS.ReadFile("icons/" + name + ".png")
	if err != nil {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	return img
}
