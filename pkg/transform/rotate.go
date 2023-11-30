package transform

import (
	"image"
	"math"
)

func RotateImage(input image.Image, angle float64) image.Image {
	bounds := input.Bounds()

	// Create a new RGBA image to hold the rotated image
	rotated := image.NewRGBA(bounds)

	// Calculate the rotation matrix
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)

	// Calculate the center of the original image
	centerX, centerY := float64(bounds.Max.X-bounds.Min.X)/2, float64(bounds.Max.Y-bounds.Min.Y)/2

	// Apply the rotation to each pixel in the original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Translate the coordinates to be relative to the center
			tx := float64(x) - centerX
			ty := float64(y) - centerY

			// Apply the rotation matrix
			newX := cosTheta*tx - sinTheta*ty
			newY := sinTheta*tx + cosTheta*ty

			// Translate the coordinates back to the original position
			newX += centerX
			newY += centerY

			// Round the new coordinates to the nearest integer
			newXInt, newYInt := int(math.Round(newX)), int(math.Round(newY))

			// Check if the new coordinates are within the bounds of the rotated image
			if newXInt >= 0 && newXInt < rotated.Bounds().Dx() && newYInt >= 0 && newYInt < rotated.Bounds().Dy() {
				rotated.Set(x, y, input.At(newXInt, newYInt))
			}
		}
	}

	return rotated
}
