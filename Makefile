ifeq ($(OS),Windows_NT)
	PLATFORM="Windows"
else
	ifeq ($(shell uname),Darwin)
		PLATFORM="MacOS"
	else
	    PLATFORM="Linux"
	endif
endif

all: test
export GO111MODULE=on
export OUTPUT=./output

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./account/... ./common/... ./xuper/...
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: test 
