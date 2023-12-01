package transform

import "image"

// Resize resizes the input image by the given factor.
// For example, A factor of 1.0 will result in no change.
// A factor of 0.5 will result in an image half the size.
func Resize(input image.Image, factor float64) image.Image {
	if factor <= 0 {
		panic("factor must be greater than 0")
	}
	bounds := input.Bounds()
	output := image.NewRGBA(image.Rect(0, 0, int(float64(bounds.Dx())*factor), int(float64(bounds.Dy())*factor)))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.Set(int(float64(x)*factor), int(float64(y)*factor), input.At(x, y))
		}
	}

	return output
}
