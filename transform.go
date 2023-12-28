package main

import (
	"image"
	"math"
)

// Rotates the input image by the specified angle in degrees.
func RotateDegrees(input image.Image, angle float64) image.Image {
	return Rotate(input, angle*math.Pi/180)
}

// Rotates the given image by the specified angle in radians.
func Rotate(img image.Image, radians float64) image.Image {
	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	// Find the new bounds
	newW, newH := rotatedDimensions(srcW, srcH, radians)
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// Calculate the center of the images
	srcCenterX, srcCenterY := float64(srcW)/2, float64(srcH)/2
	dstCenterX, dstCenterY := float64(newW)/2, float64(newH)/2

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			// Calculate the original position of this pixel
			originalX, originalY := rotatePoint(float64(x)-dstCenterX, float64(y)-dstCenterY, -radians)
			originalX += srcCenterX
			originalY += srcCenterY

			floorX, floorY := int(math.Floor(originalX)), int(math.Floor(originalY))

			if floorX >= 0 && floorX < srcW && floorY >= 0 && floorY < srcH {
				// Set the pixel from the source image if it's in bounds
				dst.Set(x, y, img.At(floorX, floorY))
			} else {
				// Set the background color, if needed
				dst.Set(x, y, COLOR.TRANSPARENT) // Transparent or black
			}
		}
	}

	return dst
}

// Rotates a point (x, y) around the origin (0, 0) by radians.
func rotatePoint(x, y, radians float64) (float64, float64) {
	sin, cos := math.Sin(radians), math.Cos(radians)
	return x*cos - y*sin, x*sin + y*cos
}

// Calculates the dimensions of a bounding box
// that fully contains a rotated w x h rectangle.
func rotatedDimensions(w, h int, radians float64) (int, int) {
	sin, cos := math.Abs(math.Sin(radians)), math.Abs(math.Cos(radians))
	newW := int(math.Round(float64(h)*sin + float64(w)*cos))
	newH := int(math.Round(float64(h)*cos + float64(w)*sin))
	return newW, newH
}

// Returns a new image that is a sub-image of the input image.
// The bounds of the new image are (x0, y0, x1, y1).
// The point (x0, y0) is the top-left corner of the new image.
// The point (x1, y1) is the bottom-right corner of the new image.
func Crop(input image.Image, x0, y0, x1, y1 int) image.Image {
	output := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			output.Set(x-x0, y-y0, input.At(x, y))
		}
	}
	return output
}

// Returns a new image that is a translated version of the input image.
// The image is translated by dx pixels horizontally and dy pixels vertically.
func Translate(input image.Image, dx, dy int) image.Image {
	bounds := input.Bounds()
	output := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x+dx < bounds.Min.X || x+dx >= bounds.Max.X || y+dy < bounds.Min.Y || y+dy >= bounds.Max.Y {
				output.Set(x+dx, y+dy, input.At(x, y))
			}
		}
	}

	return output
}
