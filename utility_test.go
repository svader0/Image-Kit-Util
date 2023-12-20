package main

import (
	"image/color"
	"os"
	"testing"
)

func TestRGBToHSL(t *testing.T) {
	// Test case 1: Red color
	h, s, l := RGBToHSL(255, 0, 0)
	if h != 0 || s != 1 || l != 0.5 {
		t.Errorf("RGBToHSL(255, 0, 0) = (%f, %f, %f), expected (0, 1, 0.5)", h, s, l)
	}

	// Test case 2: Green color
	h, s, l = RGBToHSL(0, 255, 0)
	if h != 120 || s != 1 || l != 0.5 {
		t.Errorf("RGBToHSL(0, 255, 0) = (%f, %f, %f), expected (120, 1, 0.5)", h, s, l)
	}

	// Test case 3: Blue color
	h, s, l = RGBToHSL(0, 0, 255)
	if h != 240 || s != 1 || l != 0.5 {
		t.Errorf("RGBToHSL(0, 0, 255) = (%f, %f, %f), expected (240, 1, 0.5)", h, s, l)
	}

	// Test case 4: Arbitrary color
	h, s, l = RGBToHSL(123, 221, 10)
	threshold := 0.001
	if h-87.867299 > threshold || s-0.913420 > threshold || l-0.452941 > threshold {
		t.Errorf("RGBToHSL(123, 221, 10) = (%f, %f, %f), expected (88, 0.913, 0.453)", h, s, l)
	}
}

func TestAverageRGBA(t *testing.T) {
	redColor := color.RGBA{255, 0, 0, 255}
	average := AverageRGBA(redColor)
	if average != 85 {
		t.Errorf("AverageRGBA(redColor) = %d, expected 85", average)
	}

	greyColor := color.RGBA{128, 128, 128, 255}
	average = AverageRGBA(greyColor)
	if average != 128 {
		t.Errorf("AverageRGBA(greyColor) = %d, expected 128", average)
	}

	randomColor := color.RGBA{128, 64, 192, 255}
	average = AverageRGBA(randomColor)
	if average != 128 {
		t.Errorf("AverageRGBA(randomColor) = %d, expected 128", average)
	}
}

func TestConvertToGray(t *testing.T) {
	testImage, err := LoadImage("res/test_image.jpg")
	if err != nil {
		t.Errorf("Error loading test image: %v", err)
	}

	// Convert the test image to grayscale
	grayImage := ConvertToGray(testImage)

	// Verify that the image is grayscale
	bounds := grayImage.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := grayImage.At(x, y).RGBA()
			if r != g || g != b {
				t.Errorf("ConvertToGray produced non-grayscale pixel at (%d, %d)", x, y)
			}
		}
	}
}

func TestLoadImageJPG(t *testing.T) {
	filename := "res/test_image.jpg"
	img, err := LoadImage(filename)
	if err != nil {
		t.Errorf("Error loading image: %v", err)
	}
	if img == nil {
		t.Error("LoadImage returned a nil image")
	}

	filename = "res/NOT_REAL_IMAGE_AHHHHHH__ITS_NONEXISTENT.jpg"
	img, err = LoadImage(filename)
	if err == nil {
		t.Error("LoadImage did not return an error for a nonexistent file!")
	} else {
		t.Logf("LoadImage returned error: %v", err)
	}
}

func TestCreateImage(t *testing.T) {
	// Create a red image.
	img := createImage(100, 100, COLOR.RED)

	// Verify that the image is red
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			if r != 65535 {
				t.Errorf("createImage produced non-red pixel at (%d, %d)", x, y)
			}
		}
	}

	// Verify that the image has the correct dimensions
	if img.Bounds().Max.X != 100 || img.Bounds().Max.Y != 100 {
		t.Errorf("createImage produced an image with dimensions %dx%d, expected 100x100", img.Bounds().Max.X, img.Bounds().Max.Y)
	}
}

func TestSaveImage(t *testing.T) {
	// Assuming you have an image to save, create one using createImage for testing
	testImage := createImage(100, 100, COLOR.RED)

	// Save the test image to a temporary file
	tmpFilename := "res/test_output.png"
	err := SaveImage(tmpFilename, testImage)
	// defer os.Remove(tmpFilename)
	if err != nil {
		t.Errorf("Error saving image: %v", err)
	}

	// Verify that the saved file exists
	if _, err := os.Stat(tmpFilename); os.IsNotExist(err) {
		t.Error("SaveImage did not create the output file")
	}
}

func TestSplitImage(t *testing.T) {
	testImage, err := LoadImage("res/test_image.jpg")
	if err != nil {
		t.Errorf("Error loading test image: %v", err)
	}

	// Split the test image into 3x3 grid
	images := SplitImage(testImage, 3, 3)

	// Verify that the correct number of sub-images are created
	if images.X()*images.Y() != 9 {
		t.Errorf("SplitImage produced %d sub-images, expected 9", len(images.Images))
	}
}
