#!/bin/bash

rm -f $(which networker) > /dev/null

TARGET_DIR=""
if [[ ! -z "$GOBIN" ]]; then
    TARGET_DIR=$GOBIN
fi 

if [[ ! -z "$GOPATH" ]]; then
    TARGET_DIR=$GOPATH/bin
fi

if [ -z "$TARGET_DIR" ]; then
    TARGET_DIR=/usr/local/bin
fi

if [ ! -d "$TARGET_DIR" ]; then
    mkdir -p "$TARGET_DIR"
fi

PROJECT_ROOT=$(git rev-parse --show-toplevel)
echo "installing"
    # -ldflags "-X github.com/fuskovic/networker/cmd.Version=`git describe --tags`" \
go build \
    -o $TARGET_DIR/networker \
    main.go
if [ $? -ne 0 ]; then
    echo "failed to compile networker"
    exit 1
fi

which networker
if [ $? -ne 0 ]; then
    echo "failed to install networker globally"
    exit 1
fi

networker -v
if [ $? -ne 0 ]; then
    echo "failed to identify networker version"
    exit 1
fi