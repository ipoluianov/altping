package forms

import (
	"bytes"
	"embed"
	"image"
	"image/color"
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

// disabledIcon returns a grayscale, semi-transparent copy of img
func disabledIcon(img image.Image) image.Image {
	if img == nil {
		return nil
	}
	b := img.Bounds()
	res := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			gray := uint8((299*uint32(c.R) + 587*uint32(c.G) + 114*uint32(c.B)) / 1000)
			res.SetNRGBA(x, y, color.NRGBA{R: gray, G: gray, B: gray, A: c.A * 2 / 5})
		}
	}
	return res
}
