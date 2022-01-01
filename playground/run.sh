#!/bin/bash

echo "********* GENERATING TLS CERT *************"

PROJECT_ROOT=$(git rev-parse --show-toplevel)

mkcert localhost
if [ $? -ne 0 ]; then
    echo "failed to generate TLS cert"
fi

mkdir -p $PROJECT_ROOT/playground/tls/

mv $PROJECT_ROOT/playground/localhost.pem $PROJECT_ROOT/playground/tls/localhost.pem
mv $PROJECT_ROOT/playground/localhost-key.pem $PROJECT_ROOT/playground/tls/localhost-key.pem

echo "******** SUCCESSFULLY GENERATED TLS CERT AND PRIVATE KEY ********"