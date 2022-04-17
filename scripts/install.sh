#!/bin/bash

rm -f $(which networker) > /dev/null

if [[ -z "$GOBIN" ]]; then
    echo "GOBIN unset"
    echo "add the following lines to your shells config file"
    echo "export GOPATH=\$HOME/go"
    echo "export GOBIN=\$GOPATH/bin"
    echo "export PATH=\$PATH:\$GOBIN"
    exit 1
fi

echo "installing"
PROJECT_ROOT=$(git rev-parse --show-toplevel)
go install $PROJECT_ROOT
if [ $? -ne 0 ]; then
    echo "failed to compile networker"
    exit 1
fi

networker -v
if [ $? -ne 0 ]; then
    echo "failed to validate networker installation"
    exit 1
fi