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
	go test ./account/...
	go test ./common/...
	go test ./xuper/...

.PHONY: test 
