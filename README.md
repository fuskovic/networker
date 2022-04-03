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

# Docs

* [networker list](docs/networker_list.md)	 - List information on connected network devices.
* [networker request](docs/networker_request.md)	 - Send an HTTP request.
* [networker scan](docs/networker_scan.md)	 - Scan hosts for open ports.
* [networker lookup](docs/networker_lookup.md)	 - Lookup hostnames, IPs, ISPs, nameservers, and networks.
  * [networker lookup hostname](docs/networker_lookup_hostname.md)	 - Lookup the hostname for a provided ip address.
  * [networker lookup ip](docs/networker_lookup_ip.md)	 - Lookup the ip address of the provided hostname.
  * [networker lookup isp](docs/networker_lookup_isp.md)	 - Lookup the internet service provider of a remote host.
  * [networker lookup nameservers](docs/networker_lookup_nameservers.md)	 - Lookup nameservers for the provided hostname.
  * [networker lookup network](docs/networker_lookup_network.md)	 - Lookup the network address of a provided host.