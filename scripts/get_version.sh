#!/bin/bash
PROJECT_ROOT=$(git rev-parse --show-toplevel)
git describe --tags --abbrev=0 > $PROJECT_ROOT/cmd/version.txt