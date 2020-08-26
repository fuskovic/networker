# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker/cmd/networker

# Usage 

    Usage: networker [subcommand] [flags]

    Description: A practical CLI tool for network administration.

    Commands:
            c, cap, capture  - Monitor network traffic on the LAN.
            ls, list         - List information on connected network devices.
            lu, lookup       - Lookup hostnames, IP addresses, nameservers, and networks.
            r, req, request  - Send an HTTP request.
            s, scan          - Scan the well-known ports of a given host.

# Basic Commands

## List

Useful for getting IP addresses and hostnames of devices on the LAN.

    networker ls

## Scan

Scans first 1024 ports of a given host.

    networker s --host 104.198.14.52

You can use a url for the host flag too.

    networker s --host farishuskovic.dev


## Lookup

Get the hostnames of a given address.

    networker lu --hostnames 104.198.14.52


Get the addresses of a given hostname.

    networker lookup -a farishuskovic.dev

# Advanced Commands

## Capture

Monitor network traffic on the LAN.

    networker capture

You can also use `-w` to include hostnames, sequence, and mac addresses in the output.

    networker capture -w

Write captured packets to a pcap file.

    networker capture -o capture.pcap

The pcap file specified will be created if it doesn't exist already.


## Request

Send a POST request. Optionally use JSON from a file as the body.

    networker r -u <url> -f <path> -m POST

Content-type headers are set automatically for JSON and XML files.

Add your own custom headers.

    networker r -u <url> -f <path> -m POST -a key:value,key:value


All methods are supported but if `--method` is unset, networker defaults to a GET.