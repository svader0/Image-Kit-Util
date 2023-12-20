package main

func main() {
	img, err := LoadImage("res/test_image.jpg")
	if err != nil {
		panic(err)
	}
	img = SortPixelsByHue(img)
	SaveImage("res/test_image.jpg", img)
}
