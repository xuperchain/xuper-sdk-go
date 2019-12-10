#!/bin/bash

set -e -x

cd `dirname $0`

go build -o sample ./example/sample.go


