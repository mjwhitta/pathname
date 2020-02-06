BUILD := build
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
LDFLAGS := -s -w
OUT := $(BUILD)/$(GOOS)/$(GOARCH)
SO := $(shell grep -hioPs "^package\s+\K\S+" *.go | sort -u)
SRCDIRS := $(shell find . -name "*.go" -exec dirname {} + | sort -u)

all: build

build: dir fmt lint
	@go build -ldflags "$(LDFLAGS)" -o "$(OUT)/$(SO).a"

check:
	@which go >/dev/null 2>&1

clean: fmt
	@rm -rf "$(BUILD)"

clena: clean

debug: dir fmt
	@go build -gcflags all="-l -N" -o "$(OUT)/$(SO).a"

dir:
	@mkdir -p "$(OUT)"

fmt: check
	@go fmt $(SRCDIRS) >/dev/null

gen: check
	@go generate

lint: check
	@which golint >/dev/null 2>&1 || \
	    go get -u golang.org/x/lint/golint
	@golint $(SRCDIRS)
