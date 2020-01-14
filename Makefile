all: build

build: fmt
	@go build

check:
	@which go >/dev/null 2>&1

clean: fmt

clena: clean

fmt: check
	@go fmt

gen: check
	@go generate
