#!/bin/bash

echo "purging old networker containers"

CONTAINERS=$(docker ps -aq | grep networker | awk '{ print $1 }')

if [[ ! -z "$CONTAINERS" ]]; then
    docker stop "$CONTAINERS" > /dev/null
    if [ $? -ne 0 ]; then
        echo "no containers to stop"
    fi

    docker rm "$CONTAINERS" > /dev/null
    if [ $? -ne 0 ]; then
        echo "no containers to remove"
    fi
fi

echo "deleting old networker image"
docker rmi networker
if [ $? -ne 0 ]; then
    echo "no images to delete"
fi

echo "cleaning bin"
rm -rf $(git rev-parse --show-toplevel)/bin > /dev/null
if [ $? -ne 0 ]; then
    echo "no binaries to delete"
fi
