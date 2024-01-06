package util

import (
	"image"
	"image/color"
	"sort"
)

// SortPixelsByHue sorts the pixels of the input image from left to right based on their hue and returns a new image.
func SortPixelsByHue(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Extract and store the pixel values in a slice
	pixels := make([]color.RGBA, 0, width*height)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			pixels = append(pixels, color.RGBAModel.Convert(src.At(x, y)).(color.RGBA))
		}
	}

	// Sort the pixel values based on their hue
	sort.Slice(pixels, func(i, j int) bool {
		hueI, _, _ := RGBToHSL(pixels[i].R, pixels[i].G, pixels[i].B)
		hueJ, _, _ := RGBToHSL(pixels[j].R, pixels[j].G, pixels[j].B)
		return hueI < hueJ
	})

	// Create a new image with the sorted pixel values
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			dst.Set(x, y, pixels[y*width+x])
		}
	}

	return dst
}
