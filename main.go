package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"sort"
)

func main() {
	inputImage, _ := loadImage("input2.png")

	// outputImage := transform.Crop(inputImage, 10, 10, 150, 150)
	// outputImage = Quantize(outputImage, 20)
	// outputImage := transform.RotateDegrees(inputImage, -45)
	// outputImage := Quantize(inputImage, 5, true)
	outputImage := sortPixels(inputImage)

	saveImage("output.png", outputImage)
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

func sortPixels(input image.Image) image.Image {
	bounds := input.Bounds()
	newImg := image.NewRGBA(bounds)
	pixelList := make([]color.Color, 0)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelList = append(pixelList, input.At(x, y))
		}
	}

	// sort all the pixels by their average color
	sort.Slice(pixelList, func(i, j int) bool {
		return averageRGBA(pixelList[i]) < averageRGBA(pixelList[j])
	})

	// draw the pixels in order
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, pixelList[y*bounds.Dx()+x])
		}
	}

	return newImg
}

func averageRGBA(color color.Color) uint8 {
	r, g, b, _ := color.RGBA()
	return uint8((float64(r) + float64(g) + float64(b)) / 3)
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
