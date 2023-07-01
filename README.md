# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker/v3)](https://goreportcard.com/report/github.com/fuskovic/networker/v3)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-74%25-brightgreen.svg?longCache=true&style=flat)</a>

# Features

- List devices on your LAN
- Port scanning
- Remote TTY
- DNS lookup

# Documentation

See [Docs](https://github.com/fuskovic/networker/blob/master/docs/networker.md) for command examples.

# Installation

# Download a pre-compiled binary

You can find a networker binary for your OS on the [releases](https://github.com/fuskovic/networker/releases) page.

# Install globally using Go

**Requires Go 1.18**

    go install github.com/fuskovic/networker/v3@latest

# Verify your installation

    networker -v

# Compile from source

**Requires Go 1.18**

Clone the repo, `cd` into it and run:

    make install
