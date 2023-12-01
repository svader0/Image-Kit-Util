package main

import (
	"github.com/svader0/image-kit-util/pkg/transform"
)

func main() {
	inputImage, _ := loadImage("input2.png")
	outputImage := transform.Resize(inputImage, .3)
	images := SplitImage(outputImage, 3, 3)
	for i, image := range images {
		images[i] = Quantize(image, 64, true)
	}
	outputImage = CombineImages(images, 3, 3)
	saveImage("output.png", outputImage)
}
