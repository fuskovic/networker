# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker/cmd/networker

# Usage 

    Usage: networker [subcommand] [flags]

    Description: A practical CLI tool for network administration.

    Commands:
            c, cap, capture  - Capture network packets on a given device.
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

    networker s -h 104.198.14.52

You can use a url for the host flag too.

    networker s -h farishuskovic.dev


## Lookup

Get the hostnames of a given address.

    networker lu --hostnames 104.198.14.52


Get the addresses of a given hostname.

    networker lookup -a farishuskovic.dev

# Advanced Commands

## Request

Send a POST request. Optionally use JSON from a file as the body.

    networker r -u <url> -f <path> -m POST

Content-type headers are set automatically for JSON and XML files.

Add your own custom headers.

    networker r -u <url> -f <path> -m POST -a key:value,key:value


All methods are supported but if `--method` is unset, networker defaults to a GET.

## Capture

Monitor network traffic on a device for a number of seconds.

    networker c -d en0 -s 10

Write captured packets to a pcap file.

    networker c -d en0 -s 10 -o capture.pcap

The pcap file specified will be created if it doesn't exist already.