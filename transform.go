package main

import (
	"image"
	"math"
)

// Rotates the input image by the specified angle in degrees.
func RotateDegrees(input image.Image, angle float64) image.Image {
	return Rotate(input, angle*math.Pi/180)
}

// Rotates the input image by the specified angle in radians.
func Rotate(input image.Image, angle float64) image.Image {
	bounds := input.Bounds()
	rotated := image.NewRGBA(bounds)

	// Calculate the rotation matrix
	// See: https://en.wikipedia.org/wiki/Rotation_matrix
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)

	// Find the image center point
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

			// If the new coordinates are within the bounds of the rotated image
			// then set the pixel at the new coordinates to the pixel from the
			// original image.
			if newXInt >= 0 && newXInt < rotated.Bounds().Dx() && newYInt >= 0 && newYInt < rotated.Bounds().Dy() {
				rotated.Set(x, y, input.At(newXInt, newYInt))
			}
		}
	}

	return rotated
}

func Crop(input image.Image, x0, y0, x1, y1 int) image.Image {
	output := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			output.Set(x-x0, y-y0, input.At(x, y))
		}
	}
	return output
}

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
