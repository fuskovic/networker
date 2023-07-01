#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)
rm -rf $PROJECT_ROOT/bin > /dev/null

echo "compiling for OSX"
go build -o $PROJECT_ROOT/bin/networker_darwin $PROJECT_ROOT/main.go
if [ $? -ne 0 ]; then
    echo "failed to compile networker for OSX" && exit 1
fi

echo "compiling for linux"
GOOS=linux GOARCH=amd64 go build -o $PROJECT_ROOT/bin/networker $PROJECT_ROOT/main.go
if [ $? -ne 0 ]; then
    echo "failed to compile networker for linux" && exit 1
fi

echo "compiling for windows"
GOOS=windows GOARCH=386 go build -o $PROJECT_ROOT/bin/networker.exe $PROJECT_ROOT/main.go
if [ $? -ne 0 ]; then
    echo "failed to compile networker for windows" && exit 1
fi

echo "compiled successfully"