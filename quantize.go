package util

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
)

/*
	IMPORTANT NOTE!

	The functions in this file are used to quantize an image to a specified number of colors.
	They are based on the median-cut algorithm, which is described here:
	https://en.wikipedia.org/wiki/Median_cut

	Currently, though, this implementation is terrible. It's extremely slow and horribly optimized.
	There are a lot of ways to improve it, and I'll get around to that soon, but for now just avoid using
	it for large images or large numbers of colors.
*/

// ColorBucket represents a bucket of pixels
type colorBucket struct {
	Pixels []color.RGBA
}

// ByChannel is a custom type to sort pixels based on a specific color channel
type byChannel struct {
	Pixels  []color.RGBA
	Channel int
}

func (b byChannel) Len() int      { return len(b.Pixels) }
func (b byChannel) Swap(i, j int) { b.Pixels[i], b.Pixels[j] = b.Pixels[j], b.Pixels[i] }
func (b byChannel) Less(i, j int) bool {
	return getColorChannel(b.Pixels[i], b.Channel) < getColorChannel(b.Pixels[j], b.Channel)
}

// getColorChannel returns the value of the specified color channel for a given color
func getColorChannel(c color.RGBA, channel int) uint32 {
	switch channel {
	case 0:
		return uint32(c.R)
	case 1:
		return uint32(c.G)
	case 2:
		return uint32(c.B)
	default:
		return 0
	}
}

func findMaxRange(buckets []colorBucket, channel int, resultChan chan<- struct {
	index    int
	maxRange float64
},
) {
	var maxRange float64
	var splitIndex int

	for i, bucket := range buckets {
		sort.Sort(byChannel{bucket.Pixels, channel})
		channelRange := float64(getColorChannel(bucket.Pixels[len(bucket.Pixels)-1], channel) - getColorChannel(bucket.Pixels[0], channel))
		if channelRange > maxRange {
			maxRange = channelRange
			splitIndex = i
		}
	}

	resultChan <- struct {
		index    int
		maxRange float64
	}{splitIndex, maxRange}
}

// Quantizes the image to the specified number of colors using the Median-cut algorithm.
// Returns a list of colors that represent the color palette of the image.
func GetColorPalette(img image.Image, numColors int) []color.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Put all pixels in a single bucket
	pixels := make([]color.RGBA, 0, width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixels = append(pixels, color.RGBAModel.Convert(img.At(x, y)).(color.RGBA))
		}
	}

	buckets := []colorBucket{{pixels}}

	// Repeat until the desired number of buckets is reached
	for len(buckets) < numColors {
		// Find the bucket with the greatest range in any color channel
		var maxRange float64
		var splitIndex int
		for i, bucket := range buckets {
			for channel := 0; channel < 3; channel++ {
				sort.Sort(byChannel{bucket.Pixels, channel})
				channelRange := float64(getColorChannel(bucket.Pixels[len(bucket.Pixels)-1], channel) - getColorChannel(bucket.Pixels[0], channel))
				if channelRange > maxRange {
					maxRange = channelRange
					splitIndex = i
				}
			}
		}

		// Split the chosen bucket into two
		sort.Sort(byChannel{Pixels: buckets[splitIndex].Pixels, Channel: 0})
		mid := len(buckets[splitIndex].Pixels) / 2
		newBucket := colorBucket{Pixels: buckets[splitIndex].Pixels[mid:]}
		buckets[splitIndex].Pixels = buckets[splitIndex].Pixels[:mid]

		// Add the new bucket to the list
		buckets = append(buckets, newBucket)
	}

	// Average the pixels in each bucket to get the final color palette
	palette := make([]color.RGBA, 0, numColors)
	for _, bucket := range buckets {
		avgColor := averageColor(bucket.Pixels)
		palette = append(palette, avgColor)
	}

	return palette
}

// Calculates the average color of a set of pixels
func averageColor(pixels []color.RGBA) color.RGBA {
	var sumR, sumG, sumB, sumA uint32
	for _, p := range pixels {
		sumR += uint32(p.R)
		sumG += uint32(p.G)
		sumB += uint32(p.B)
		sumA += uint32(p.A)
	}

	size := uint32(len(pixels))
	return color.RGBA{
		R: uint8(sumR / size),
		G: uint8(sumG / size),
		B: uint8(sumB / size),
		A: uint8(sumA / size),
	}
}

// Quantizes the input image using the median-cut algorithm
func Quantize(img image.Image, numColors int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	palette := GetColorPalette(img, numColors)

	// Create a new RGBA image to store the result
	result := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over each pixel in the original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the original pixel color
			originalColor := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			// Find the nearest color in the palette
			nearestColor := findNearestColor(originalColor, palette)

			// Set the new color in the result image
			result.Set(x, y, nearestColor)
		}
	}

	return result
}

// Finds the nearest color in the palette to the given color
func findNearestColor(target color.RGBA, palette []color.RGBA) color.RGBA {
	var minDist uint32 = 0xFFFFFFFF
	var nearestColor color.RGBA

	for _, c := range palette {
		// Calculate Euclidean distance in RGB space
		dist := colorDistanceSquared(target, c)
		if dist < minDist {
			minDist = dist
			nearestColor = c
		}
	}

	return nearestColor
}

// Calculates the squared Euclidean distance between two colors in RGB space
func colorDistanceSquared(c1, c2 color.RGBA) uint32 {
	dr := uint32(c1.R) - uint32(c2.R)
	dg := uint32(c1.G) - uint32(c2.G)
	db := uint32(c1.B) - uint32(c2.B)
	da := uint32(c1.A) - uint32(c2.A)

	return dr*dr + dg*dg + db*db + da*da
}

/*
IMPORTANT NOTE!
QuantizeKMeans is TERRIBLY OPTIMIZED, and NOT PERFORMANT AT ALL!!
Do not use it for anything of actual importance, as it will likely take ages to run.
However, the end result does look pretty cool, so I'm keeping it here for now.
*/
func QuantizeKMeans(input image.Image, numColors int, dither bool) image.Image {
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
	for iter := 0; iter < 8; iter++ {
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
