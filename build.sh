#!/bin/bash

set -e -x

cd `dirname $0`

#go build -o main ./example/main.go
go build -o sample ./example/sample.go

# build plugins
echo "OS:"${PLATFORM}
echo "## Build Plugins..."
mkdir -p plugins/crypto
go build --buildmode=plugin -o plugins/crypto/crypto-default.so.1.0.0 ./crypto/client/xchain/xchain_crypto_client.go


