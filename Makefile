all: build

build: fmt
	@go build .

fmt:
	@go fmt .
