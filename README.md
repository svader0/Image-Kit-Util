# image-kit-util (WORK IN PROGRESS)

Image-Kit-Util provides a simple and clean interface for image manipulation in Go. Perfect for generative art and automated image manipulation scripts.

### Supported operations:
- Crop
- Resize
- Rotate
- Translate
- Sort pixels by hue
- Quantize (w/ Floyd-Steinberg Dithering)
- Greyscale conversion
- Split/Combine images

## Example

```Go
package main

import (
	"github.com/svader0/image-kit-util/pkg/transform"
)

func main() {
    // First, load our image from the local directory.
	inputImage, err := loadImage("./input.jpg")
    if err != nil {
        panic("Could not load image!")
    }

    // Split our image into 9 smaller images
	images := SplitImage(outputImage, 3, 3)
	for i, image := range images {
        // Quantize all our images with dithering enabled
		images[i] = Quantize(image, 15, true)
	}
    // Recombine the 9 small images into a single image.
	outputImage = CombineImages(images, 3, 3)

    // Save our output as a .png file.
	saveImage("output.png", outputImage)
}

```

## TODO / Needs Work

- Current quantize and pixelsort algorithms are SLOW! Images over 1000x1000 take significant time to render.
- Functions aren't organized. Working on this.
- New features to add soon:
    - Drawing shapes, lines, text, etc. on image.
    - Support for image overlays
    - Digital Solarization
    - Hue shift
    - Range hue shift