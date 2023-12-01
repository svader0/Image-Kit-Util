package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// FOR COLOR IMAGES -------------------------------------------------

// Quantize performs color quantization on the input image.
// The numColors parameter specifies the number of colors to quantize to.
// The dither parameter specifies whether to apply Floyd-Steinberg dithering.
func Quantize(input image.Image, numColors int, dither bool) image.Image {
	palette := make(color.Palette, 0, numColors)
	bounds := input.Bounds()

	// Initialize the palette with unique colors from the image
	colorSet := make(map[color.Color]bool)
	for y := bounds.Min.Y; y < bounds.Max.Y && len(colorSet) < numColors; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && len(colorSet) < numColors; x++ {
			colorSet[input.At(x, y)] = true
		}
	}

	// Convert the color set to a palette
	for color := range colorSet {
		palette = append(palette, color)
	}

	// Perform k-means clustering
	for iter := 0; iter < 25; iter++ {
		// Assign each pixel to the nearest color in the palette
		assignments := make([]int, numColors)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				pixel := input.At(x, y)
				nearestIndex := findNearestColorIndex(pixel, palette)
				assignments[nearestIndex]++
			}
		}

		// Update the palette by computing the mean color of each cluster
		for i := 0; i < numColors; i++ {
			sumR, sumG, sumB, count := 0, 0, 0, assignments[i]
			if count > 0 {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						pixel := input.At(x, y)
						nearestIndex := findNearestColorIndex(pixel, palette)
						if nearestIndex == i {
							r, g, b, _ := pixel.RGBA()
							sumR += int(r)
							sumG += int(g)
							sumB += int(b)
						}
					}
				}

				palette[i] = color.RGBA{
					uint8(sumR / count),
					uint8(sumG / count),
					uint8(sumB / count),
					255,
				}
			}
		}
	}
	output := image.NewPaletted(bounds, palette)
	if dither {
		draw.FloydSteinberg.Draw(output, bounds, input, image.Point{})
	} else {
		draw.Draw(output, bounds, input, image.Point{}, draw.Src)
	}
	return convertPalettedToImage(output)
}

// ConvertPalettedToImage converts a paletted image to a regular image.Image.
func convertPalettedToImage(paletted *image.Paletted) image.Image {
	bounds := paletted.Bounds()
	rgba := image.NewRGBA(bounds)

	// Create a color palette from the paletted image
	palette := color.Palette(paletted.Palette)

	// Draw the paletted image onto the RGBA image
	draw.Draw(rgba, bounds, paletted, bounds.Min, draw.Over)

	// Replace the color indices in the RGBA image with actual colors from the palette
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, palette.Convert(rgba.At(x, y)))
		}
	}

	return rgba
}

// findNearestColorIndex finds the index of the nearest color in the palette.
func findNearestColorIndex(c color.Color, palette color.Palette) int {
	cr, cg, cb, _ := c.RGBA()
	minDist := math.MaxUint32
	nearestIndex := 0

	for i, p := range palette {
		pr, pg, pb, _ := p.RGBA()
		dist := sqrDiff(cr, pr) + sqrDiff(cg, pg) + sqrDiff(cb, pb)
		if dist < uint32(minDist) {
			minDist = int(dist)
			nearestIndex = i
		}
	}

	return nearestIndex
}

// sqrDiff computes the squared difference between two values.
func sqrDiff(a, b uint32) uint32 {
	if a > b {
		return (a - b) * (a - b)
	}
	return (b - a) * (b - a)
}

// findNearestLevelIndex finds the index of the nearest level.
func findNearestLevelIndex(value uint8, levels []int) int {
	var minDist uint32 = math.MaxUint32
	nearestIndex := 0

	for i, level := range levels {
		dist := sqrDiff(uint32(value), uint32(level))
		if dist < minDist {
			minDist = dist
			nearestIndex = i
		}
	}

	return nearestIndex
}
