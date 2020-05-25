#!/bin/bash
protoc -I. xchain-spv.proto --go_out=plugins=grpc:. -I ./../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/
protoc -I. xchain.proto --go_out=plugins=grpc:. -I ./../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/
protoc -I. chainedbft.proto --go_out=plugins=grpc:. -I ./../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/
protoc -I. xendorser.proto --go_out=plugins=grpc:. -I ./../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/