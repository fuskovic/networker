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

# Commands

## List

Useful for getting IP addresses and hostnames of devices on the LAN.

    networker ls

## Scan

Scans first 1024 ports of a given host.

    networker s --host 104.198.14.52

You can use a url for the host flag too.

    networker s --host farishuskovic.dev

If you forget to provide an http or https proto-scheme in your URL networker will append it for you.


## Lookup

Get the hostnames of a given address.

    networker lu --hostnames 104.198.14.52


Get the addresses of a given hostname.

    networker lu --addresses farishuskovic.dev


## Capture

Monitor network traffic on the LAN.

    networker cap

You can save your capture session to a pcap file.

    networker cap -o capture.pcap

If the file doesn't exist, networker will create it for you.

You can also use `--wide` to include hostnames, sequence-numbers, and mac addresses in the output.

    networker cap --wide


## Request

Here's an example of how to send a post request using a JSON file as the request body.

    networker req -m POST -f /path/to/file.json -u https://url.com

Content-type headers are automatically set when JSON and XML files are provided.

You can add your own custom headers as a comma-separated list of key/value pairs.

    networker req --add-headers key:value,key:value -u https://url.com

All methods are supported but if `--method` is unset, networker will default to a GET.

