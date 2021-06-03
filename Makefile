ifeq ($(OS),Windows_NT)
	PLATFORM="Windows"
else
	ifeq ($(shell uname),Darwin)
		PLATFORM="MacOS"
	else
	    PLATFORM="Linux"
	endif
endif

all: build
export GO111MODULE=on
export GOFLAGS=-mod=vendor
export OUTPUT=./output

build:
	PLATFORM=$(PLATFORM) ./build.sh

test:
	go test ./account/...
	go test ./common/...
	go test ./xuper/...

clean:
	rm -rf main
	rm -rf sample

.PHONY: all test clean
