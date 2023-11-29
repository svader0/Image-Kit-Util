package main

import (
	"image"
	"image/color"
)

// floydSteinbergDithering applies the Floyd-Steinberg dithering algorithm to a grayscale image
func floydSteinbergDithering(input *image.Gray) *image.Gray {
	bounds := input.Bounds()
	output := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := input.GrayAt(x, y)
			newPixel := color.GrayModel.Convert(oldPixel).(color.Gray)
			output.SetGray(x, y, newPixel)

			quantError := int(oldPixel.Y - newPixel.Y)

			// Distribute the error to neighboring pixels
			if x+1 < bounds.Max.X {
				distributeError(input, x+1, y, quantError, 7)
			}

			if y+1 < bounds.Max.Y {
				if x-1 > bounds.Min.X {
					distributeError(input, x-1, y+1, quantError, 3)
				}

				distributeError(input, x, y+1, quantError, 5)

				if x+1 < bounds.Max.X {
					distributeError(input, x+1, y+1, quantError, 1)
				}
			}
		}
	}

	return output
}

// distributeError distributes quantization error to a neighboring pixel
func distributeError(input *image.Gray, x, y, quantError, weight int) {
	oldPixel := input.GrayAt(x, y)
	newPixel := color.Gray{Y: uint8(int(oldPixel.Y) + quantError*weight/16)}
	input.SetGray(x, y, newPixel)
}

// floydSteinbergDitheringColor applies the Floyd-Steinberg dithering algorithm to a color image
func floydSteinbergDitheringColor(input image.Image) image.Image {
	bounds := input.Bounds()
	output := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := input.At(x, y)
			newPixel := color.RGBAModel.Convert(oldPixel).(color.RGBA)
			output.SetRGBA(x, y, newPixel)

			quantErrorR := int(oldPixel.(color.RGBA).R - newPixel.R)
			quantErrorG := int(oldPixel.(color.RGBA).G - newPixel.G)
			quantErrorB := int(oldPixel.(color.RGBA).B - newPixel.B)

			// Distribute the error to neighboring pixels
			if x+1 < bounds.Max.X {
				distributeErrorColor(output, x+1, y, quantErrorR, quantErrorG, quantErrorB, 7)
			}
			if y+1 < bounds.Max.Y {
				if x-1 > bounds.Min.X {
					distributeErrorColor(output, x-1, y+1, quantErrorR, quantErrorG, quantErrorB, 3)
				}
				distributeErrorColor(output, x, y+1, quantErrorR, quantErrorG, quantErrorB, 5)
				if x+1 < bounds.Max.X {
					distributeErrorColor(output, x+1, y+1, quantErrorR, quantErrorG, quantErrorB, 1)
				}
			}
		}
	}

	return output
}

// distributeErrorColor distributes quantization error to a neighboring pixel in a color image
func distributeErrorColor(input *image.RGBA, x, y, quantErrorR, quantErrorG, quantErrorB, weight int) {
	oldPixel := input.RGBAAt(x, y)
	newPixel := color.RGBA{
		R: uint8(int(oldPixel.R) + quantErrorR*weight/16),
		G: uint8(int(oldPixel.G) + quantErrorG*weight/16),
		B: uint8(int(oldPixel.B) + quantErrorB*weight/16),
		A: oldPixel.A,
	}
	input.SetRGBA(x, y, newPixel)
}
