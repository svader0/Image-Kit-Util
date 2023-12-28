name := image-kit-util

.PHONY: build run clean test

build:
	go build -o bin/$(name) main.go utility.go quantize.go fun.go transform.go

run: build
	./bin/$(name)

clean:
	go clean
	if [ -f bin/$(name) ] ; then rm -f -r bin/$(name) ; fi
	find ./res -type f ! \( -name 'test_image.jpg' -o -name 'test_output.jpg' \) -exec rm {} +

test: build
	go test 