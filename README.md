# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker

# Commands

## List

    networker list

<img src="gifs/list.gif" height="400" width="1300">

## Lookup

    networker lookup --hostnames 104.198.14.52
    networker lookup --network farishuskovic.dev
    networker lookup --nameservers farishuskovic.dev
    networker lookup --addresses farishuskovic.dev


<img src="gifs/lookup.gif" height="400" width="1300">

## Scan

    networker scan --host farishuskovic.dev -v

<img src="gifs/scan.gif" height="100" width="1000">


## Request

    networker request --url https://api.thecatapi.com/v1/breeds

<img src="gifs/request.gif" height="200" widht="1000">


## Capture

    networker capture -d en0 -s 10

<img src="gifs/cap.gif" height="400" width="8130">

