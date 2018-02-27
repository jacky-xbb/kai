help:
	@echo '    build ...................... generates a binary for kai server'
	@echo '    run ........................ runs kai server'
	@echo '    install .................... install kai server'
	@echo '    uninstall .................. uninstall kai server'
	@echo '    clean ...................... clean files have generates'
	@echo '    check ...................... runs golint and vet tools'
	@echo '    test ....................... runs tests locally'
	@echo '    test_coverage .............. runs tests and generates coverage profile'

SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)
COVERAGE_FILE := coverage.txt 

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run update

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET) $(COVERAGE_FILE)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)

test: build
	@ginkgo -r --slowSpecThreshold=20 --succinct .

test_coverage: build
	@ginkgo -r --slowSpecThreshold=20 --cover --succinct .
	@gover
	@mv gover.coverprofile $(COVERAGE_FILE)

# Glide update
update:
	@glide u