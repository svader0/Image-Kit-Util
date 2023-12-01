package transform

import (
	"image"
)

func Crop(input image.Image, x0, y0, x1, y1 int) image.Image {
	output := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			output.Set(x-x0, y-y0, input.At(x, y))
		}
	}
	return output
}
