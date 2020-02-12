#!/bin/bash

set -e -x

cd `dirname $0`

go build -o sample ./example/sample.go
go build example/main_counter.go
go build example/main_trust_counter.go
