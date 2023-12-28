package main

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

// Color constants
var COLOR = struct {
	RED         color.RGBA
	BLUE        color.RGBA
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
	BLUE:        color.RGBA{0, 0, 255, 255},
	BLACK:       color.RGBA{0, 0, 0, 255},
	WHITE:       color.RGBA{255, 255, 255, 255},
	YELLOW:      color.RGBA{255, 255, 0, 255},
	ORANGE:      color.RGBA{255, 165, 0, 255},
	GREEN:       color.RGBA{0, 255, 0, 255},
	PURPLE:      color.RGBA{128, 0, 128, 255},
	PINK:        color.RGBA{255, 192, 203, 255},
	TRANSPARENT: color.RGBA{0, 0, 0, 0},
}

// Converts RGB values to HSL.
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
		h = 60 * math.Mod(((gF-bF)/delta)+6, 6)
		if h < 0 {
			h += 360
		}
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

	return h, s, l
}

// Finds the average of the RGB values of a color
func AverageRGBA(color color.Color) uint8 {
	r, g, b, _ := color.RGBA()
	return uint8((uint32(r) + uint32(g) + uint32(b)) / 3)
}

// Converts an image to grayscale
func ConvertToGray(input image.Image) *image.Gray {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := input.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	return gray
}

// Loads an image from file
func LoadImage(filename string) (image.Image, error) {
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

// Saves an image to file
func SaveImage(filename string, img image.Image) error {
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

// Represents a group of subimages that can be combined into a single image
// The subimages are arranged in a grid of size x by y.
type Subimages struct {
	Images []image.Image
	x      int
	y      int
}

// X returns the x value.
func (s *Subimages) X() int {
	return s.x
}

// Y returns the y value.
func (s *Subimages) Y() int {
	return s.y
}

// Splits the image into a grid of smaller images.
// Returns a Subimages struct containing the smaller images and the arrangement of the grid. (x by y)
func SplitImage(img image.Image, x, y int) Subimages {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	images := make([]image.Image, 0, x*y)

	// Calculate the width and height of each sub-image
	subWidth := width / x
	subHeight := height / y

	// Create a new image for each sub-image
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			subImage := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(i*subWidth, j*subHeight, (i+1)*subWidth, (j+1)*subHeight))
			images = append(images, subImage)
		}
	}

	return Subimages{images, x, y}
}

// Takes in a subimages struct and combines them into a single image
func CombineImages(imgs *Subimages) image.Image {
	bounds := imgs.Images[0].Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a new image for the result
	combinedImage := image.NewRGBA(image.Rect(0, 0, width*imgs.X(), height*imgs.Y()))

	// Combine the images
	for i, img := range imgs.Images {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				combinedImage.Set(i*width+x, i/width*height+y, img.At(x, y))
			}
		}
	}

	return combinedImage
}

func AverageColor(img *image.Image) color.RGBA {
	bounds := (*img).Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var r, g, b uint32
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rTemp, gTemp, bTemp, _ := (*img).At(x, y).RGBA()
			r += rTemp
			g += gTemp
			b += bTemp
		}
	}

	r /= uint32(width * height)
	g /= uint32(width * height)
	b /= uint32(width * height)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
