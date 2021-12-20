#!/bin/bash

echo "********* GENERATING TLS CERT *************"

PROJECT_ROOT=$(git rev-parse --show-toplevel)

go run \
    /usr/local/go/src/crypto/tls/generate_cert.go \
    --rsa-bits 1024 \
    --host 127.0.0.1,::1,localhost \
    --ca --start-date "Jan 1 00:00:00 1970" \
    --duration=1000000h

if [ $? -ne 0 ]; then
    echo "failed to generate TLS cert"
fi

mv $(PROJECT_ROOT)/scripts/cert.pem $(PROJECT_ROOT)/internal/test/cert.pem
mv $(PROJECT_ROOT)/scripts/key.pem test/key.pem

echo "******** SUCCESSFULLY GENERATED TLS CERT ********"