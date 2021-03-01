# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker/cmd/networker

# Usage 

    Usage: networker [subcommand] [flags]

    Description: A practical CLI tool for network administration.

    Commands:
            ls, list         - List information on connected network devices.
            lu, lookup       - Lookup hostnames, IP addresses, nameservers, and networks.
            r, req, request  - Send an HTTP request.
            s, scan          - Scan the well-known ports of a given host.

# Commands

## List

Useful for getting IP addresses and hostnames of devices on the LAN.

    networker ls

## Scan

Scans the well-known ports of a given host:

    networker scan --host <ip-address>

Scan all ports of a given host:

    networker scan --host <ip-address> --all

Hostnames are also supported if you don't wan't to use an ip address:

    networker s --host <hostname>

If you don't use the `--host` flag then the `scan` command will scan all devices on your local network:

    networker scan


## Lookup

Get the hostname of an IP address:

    networker lookup hostname --ip <ip-address>


Get the ip address of a hostname:

    networker lookup ip --hostname <hostname>



## Request

The `request` command defaults to a `GET` if the `--method` flag isn't passed:

    networker request --url <url-endpoint>

You can provision the request body from a `json` or `xml` file when sending `POST` requests:

    networker request --method POST --file <file-path> --url <url-endpoint>

If you pass a `json` or `xml` file the `Content-Type` header is automatically set for you.

If you wan't to add your own custom headers as a comma-separated list of key/value pairs.

    networker request --add-headers key:value,key:value --url <url-endpoint>

You don't need to use quotes when adding custom headers.

