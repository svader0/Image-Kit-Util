package main

func main() {
	img, err := LoadImage("res/test_image.jpg")
	if err != nil {
		panic(err)
	}

	// Crop the image to a square
	img = Crop(img, 0, 0, 500, 500)

	// Rotate the image by 45 degrees
	img = RotateDegrees(img, 45)

	// Convert the image to grayscale
	img = ConvertToGray(img)

	// Save the image
	err = SaveImage("res/test_output.jpg", img)
	if err != nil {
		panic(err)
	}
}
