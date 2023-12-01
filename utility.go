package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Color constants
var COLOR = struct {
	RED         color.RGBA
	BLACK       color.RGBA
	WHITE       color.RGBA
	YELLOW      color.RGBA
	ORANGE      color.RGBA
	GREEN       color.RGBA
	PURPLE      color.RGBA
	PINK        color.RGBA
	TRANSPARENT color.RGBA
}{
	RED:         color.RGBA{255, 0, 0, 255},
	BLACK:       color.RGBA{0, 0, 0, 255},
	WHITE:       color.RGBA{255, 255, 255, 255},
	YELLOW:      color.RGBA{255, 255, 0, 255},
	ORANGE:      color.RGBA{255, 165, 0, 255},
	GREEN:       color.RGBA{0, 255, 0, 255},
	PURPLE:      color.RGBA{128, 0, 128, 255},
	PINK:        color.RGBA{255, 192, 203, 255},
	TRANSPARENT: color.RGBA{0, 0, 0, 0},
}

// RGBToHSL converts RGB values to HSL.
func RGBToHSL(r, g, b uint8) (float64, float64, float64) {
	rF, gF, bF := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0

	max := math.Max(math.Max(rF, gF), bF)
	min := math.Min(math.Min(rF, gF), bF)
	delta := max - min

	var h, s, l float64

	// Calculate hue
	if delta == 0 {
		h = 0
	} else if max == rF {
		h = 60 * math.Mod(((gF-bF)/delta), 6)
	} else if max == gF {
		h = 60 * (((bF - rF) / delta) + 2)
	} else {
		h = 60 * (((rF - gF) / delta) + 4)
	}

	// Calculate lightness
	l = (max + min) / 2

	// Calculate saturation
	if delta == 0 {
		s = 0
	} else {
		s = delta / (1 - math.Abs(2*l-1))
	}

	return math.Mod(h+360, 360), s, l
}

// Finds the average of the RGB values of a color
func averageRGBA(color color.Color) uint8 {
	r, g, b, _ := color.RGBA()
	return uint8((float64(r) + float64(g) + float64(b)) / 3)
}

// Converts an image to grayscale
func ConvertToGray(input image.Image) *image.Gray {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, input.At(x, y))
		}
	}

	return gray
}

// loadImage loads an image from file
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func createImage(x, y int, color color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, x, y))
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			img.Set(i, j, color)
		}
	}
	return img
}

// saveImage saves an image to file
func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// Splits the image into a grid of smaller images
func SplitImage(img image.Image, x, y int) []image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	images := make([]image.Image, 0, x*y)

	// Calculate the width and height of each sub-image
	subWidth := width / x
	subHeight := height / y

	// Create a new image for each sub-image
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			subImage := image.NewRGBA(image.Rect(0, 0, subWidth, subHeight))
			for k := 0; k < subWidth; k++ {
				for l := 0; l < subHeight; l++ {
					subImage.Set(k, l, img.At(i*subWidth+k, j*subHeight+l))
				}
			}
			images = append(images, subImage)
		}
	}

	return images
}

// Takes in a slice of subimages and combines them into a single image
// The subimages are arranged in a grid of size x by y.
// This is meant to be used with the output of SplitImage, but it can be used with any slice of images.
func CombineImages(images []image.Image, x, y int) image.Image {
	bounds := images[0].Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a new image for the combined image
	combinedImage := image.NewRGBA(image.Rect(0, 0, width*x, height*y))

	// Combine the images
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < width; k++ {
				for l := 0; l < height; l++ {
					combinedImage.Set(i*width+k, j*height+l, images[i*x+j].At(k, l))
				}
			}
		}
	}

	return combinedImage
}
