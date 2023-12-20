name := image-kit-util

.PHONY: build run clean test

build:
	go build -o bin/$(name) main.go utility.go quantize.go fun.go translate.go

run: build
	./bin/$(name)

clean:
	go clean
	if exists bin\$(name) del bin\$(name)
	if exists bin\$(name).exe del bin\$(name).exe
	if exists bin del bin

test: build
	go test 