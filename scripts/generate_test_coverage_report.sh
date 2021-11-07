#!/bin/bash

echo "***** GENERATING COVERAGE REPORT *****"

go test -coverprofile coverage.out ./...

if [ $? -ne 0 ]; then
    echo "tests failing"
    echo "coverage report not generated"
    exit 1
fi

if [ "$1" = "headless" ]; then
    echo "***** CHECKING COVERAGE *****"
    go tool cover -func coverage.out
else
    echo "***** SERVING HTML COVERAGE DIFF... *****"
    go tool cover -html coverage.out
fi

rm coverage.out