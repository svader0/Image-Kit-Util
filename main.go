package main

import (
	"image"
	"image/png"
	"os"
)

func main() {
	// Load the input image
	inputImage, err := loadImage("input3.png")
	if err != nil {
		panic(err)
	}

	// Convert the input image to grayscale
	// grayImage := convertToGray(inputImage)

	// Apply Floyd-Steinberg dithering
	// ditheredImage := floydSteinbergDithering(grayImage)
	quantizedImage := QuantizeImage(inputImage, 8)

	// Save the dithered image
	err = saveImage("output.png", quantizedImage)
	if err != nil {
		panic(err)
	}
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

// convertToGray converts an image to grayscale
func convertToGray(input image.Image) *image.Gray {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, input.At(x, y))
		}
	}

	return gray
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
