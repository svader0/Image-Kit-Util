package transform

import "image"

func Translate(input image.Image, dx, dy int) image.Image {
	bounds := input.Bounds()
	output := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.Set(x+dx, y+dy, input.At(x, y))
		}
	}

	return output
}
