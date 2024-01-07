package util

import (
	"image"
	"image/color"
	"image/draw"
)

// Adds a border to the given image.
func AddBorder(img image.Image, size int, borderColor color.Color) *image.RGBA {
	// Calculate new dimensions.
	newWidth := img.Bounds().Dx() + (size * 2)
	newHeight := img.Bounds().Dy() + (size * 2)

	// Create a new image with the new dimensions.
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Set the border area to the specified color.
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{borderColor}, image.Point{}, draw.Src)

	// Draw original image on top of the border.
	draw.Draw(newImg, image.Rect(size, size, size+img.Bounds().Dx(), size+img.Bounds().Dy()), img, img.Bounds().Min, draw.Over)

	return newImg
}
