package transform

import (
	"image"
	"image/draw"
)

func Crop(input image.Image, x0, y0, x1, y1 int) image.Image {
	// bounds := input.Bounds()
	output := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	draw.Draw(output, output.Bounds(), input, image.Point{x0, y0}, draw.Src)
	return output
}
