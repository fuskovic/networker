# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker/cmd/networker

# Usage 

    Usage: networker [subcommand] [flags]

    A practical CLI tool for network administration.

    Commands:
            capture  Capture network packets on a given device.
            list     List information on connected network devices.
            lookup   Lookup hostnames, IP addresses, nameservers, and networks.
            request  Send an HTTP request.
            scan     Scan the well-known ports of a given host.

# Examples

## List

    networker list

<img src="gifs/list.gif" height="400" width="1300">

## Scan

    networker scan --host farishuskovic.dev -v

<img src="gifs/scan.gif" height="100" width="1000">


## Lookup

    networker lookup --hostnames 104.198.14.52
    networker lookup --network farishuskovic.dev
    networker lookup --nameservers farishuskovic.dev
    networker lookup --addresses farishuskovic.dev


<img src="gifs/lookup.gif" height="400" width="1300">


