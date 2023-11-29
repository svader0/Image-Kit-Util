name := image-kit-util

build:
	go build -o bin/$(name) main.go dithering.go quantize.go

run: build
	./bin/$(name)

clean:
	rm -rf bin