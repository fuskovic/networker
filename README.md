# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-56%25-brightgreen.svg?longCache=true&style=flat)</a>


# Installation

## Install by downloading pre-compiled binaries

Checkout the [releases](https://github.com/fuskovic/networker/releases) page to download the latest executables for Linux, Mac, and Windows.

## Install globally using Go

    go install github.com/fuskovic/networker@latest

Then verify your installation:

    networker -v

## Install by compiling from source

Clone the repo, `cd` into it and run:

    make install

# Usage 

```
Usage: networker [subcommand] [flags]

Description: A simple networking tool.

networker flags:
  -v, --version   Installed version.

Commands:
	ls, list         - List information on connected network devices.
	lu, lookup       - Lookup hostnames, IP addresses, internet service providers, nameservers, and networks.
	r, req, request  - Send an HTTP request.
	s, scan          - Scan hosts for open ports.
```

# Commands 

## List 

```
Usage: networker list [flags]

Aliases: ls

Description: List information on connected network devices.

networker list flags:
      --json   Output as json.
```

## Lookup 

```
Usage: networker lookup [subcommand]

Aliases: lu

Description: Lookup hostnames, IP addresses, internet service providers, nameservers, and networks.

Commands:
	hostname     - Lookup the hostname for a provided ip address.
	ip           - Lookup the ip address of the provided hostname.
	network      - Lookup the network address of a provided host.
	nameservers  - Lookup nameservers for the provided hostname.
	isp          - Lookup the internet service provider of a remote host.
```

## Scan 

```
Usage: networker scan [flags] [host]

Aliases: s

Description: Scan hosts for open ports.

networker scan flags:
  -a, --all    Scan all ports(scans first 1024 if not enabled).
      --json   Output as json.
```

## Request 

```
Usage: networker request [flags] [url]

Aliases: r, req

Description: Send an HTTP request.

networker request flags:
  -b, --body string       Request body. (you can use a JSON string literal or a path to a json file)
  -f, --files string      Files to upload. (format: formname=path/to/file1,path/to/file2,path/to/file3)
  -H, --headers strings   Request headers.(format: key:value,key:value,key:value)
  -j, --json-only         Only output json.
  -m, --method string     Request method. (default "GET")
```
