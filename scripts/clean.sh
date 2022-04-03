#!/bin/bash

echo "cleaning bin"
rm -rf $(git rev-parse --show-toplevel)/bin > /dev/null
if [ $? -ne 0 ]; then
    echo "no binaries to delete"
fi
