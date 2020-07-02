#!/bin/bash

go build -o $OUTPUT/main ./example/main.go

go build -o $OUTPUT/gm_main ./example/gm_main.go

go build -o $OUTPUT/sample ./example/sample.go

cp -r ./conf $OUTPUT/