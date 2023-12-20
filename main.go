package main

func main() {
	img, err := LoadImage("res/test_image.jpg")
	if err != nil {
		panic(err)
	}

	// Rotate the image by 45 degrees
	img = RotateDegrees(img, 45)

	// Crop the image to a square
	img = Crop(img, 0, 0, 500, 500)

	// Convert the image to grayscale
	img = ConvertToGray(img)

	// Save the image
	err = SaveImage("out.jpg", img)
	if err != nil {
		panic(err)
	}
}
