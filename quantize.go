package main

import (
	"image"
	"image/color"
	"image/draw"
)

// QuantizeImage quantizes the input image using k-means algorithm with the specified number of colors.
func QuantizeImage(inputImage image.Image, numColors int) image.Image {
	// Convert the image to the RGBA color model
	bounds := inputImage.Bounds()
	rgbaImage := image.NewRGBA(bounds)
	draw.Draw(rgbaImage, bounds, inputImage, bounds.Min, draw.Src)

	// Flatten the image into a 2D slice of color.RGBA values
	pixels := make([]color.RGBA, 0, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, rgbaImage.At(x, y).(color.RGBA))
		}
	}

	// Perform k-means clustering to find dominant colors
	centroids := kMeans(pixels, numColors)

	// Replace each pixel with the nearest centroid color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			index := nearestCentroidIndex(rgbaImage.At(x, y).(color.RGBA), centroids)
			rgbaImage.Set(x, y, centroids[index])
		}
	}

	return rgbaImage
}

// kMeans performs k-means clustering on the input colors to find dominant colors.
func kMeans(colors []color.RGBA, k int) []color.RGBA {
	//
}

// nearestCentroidIndex finds the index of the nearest centroid for a given color.
func nearestCentroidIndex(c color.RGBA, centroids []color.RGBA) int {
	minDist := colorDistance(c, centroids[0])
	index := 0

	for i := 1; i < len(centroids); i++ {
		dist := colorDistance(c, centroids[i])
		if dist < minDist {
			minDist = dist
			index = i
		}
	}

	return index
}

// colorDistance calculates the Euclidean distance between two colors.
func colorDistance(c1, c2 color.RGBA) float64 {
	r := float64(c1.R) - float64(c2.R)
	g := float64(c1.G) - float64(c2.G)
	b := float64(c1.B) - float64(c2.B)
	return r*r + g*g + b*b
}
