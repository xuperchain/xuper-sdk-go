#!/bin/bash
cp -r ./conf $OUTPUT/

go build -o $OUTPUT/main ./example/main.go

go build -o $OUTPUT/sample ./example/sample.go

