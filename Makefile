name := image-kit-util

build:
	go build -o bin/$(name) main.go quantize.go 

run: build
	./bin/$(name)

clean:
	rm bin/$(name)