#!/bin/bash

echo "********* GENERATING TLS CERT *************"

PROJECT_ROOT=$(git rev-parse --show-toplevel)

if [[ -z "$GOROOT" ]]; then
    echo "GOROOT environment variable not set"
    exit 1
fi

go run \
    $GOROOT/src/crypto/tls/generate_cert.go \
    --rsa-bits 1024 \
    --host 127.0.0.1,::1,localhost \
    --ca --start-date "Jan 1 00:00:00 1970" \
    --duration=1000000h

if [ $? -ne 0 ]; then
    echo "failed to generate TLS cert"
fi

mkdir -p $PROJECT_ROOT/playground/tls/

mv $PROJECT_ROOT/playground/cert.pem $PROJECT_ROOT/playground/tls/cert.pem
mv $PROJECT_ROOT/playground/key.pem $PROJECT_ROOT/playground/tls/key.pem

echo "******** SUCCESSFULLY GENERATED TLS CERT AND PRIVATE KEY ********"