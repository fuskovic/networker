# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker

# Commands

## List

    networker list

![list](gifs/list.gif)

## Lookup

    networker lookup --hostnames 104.198.14.52
    networker lookup --network farishuskovic.dev
    networker lookup --nameservers farishuskovic.dev
    networker lookup --addresses farishuskovic.dev

![lookup](gifs/lookup.gif)

## Request

    networker request --url https://api.thecatapi.com/v1/breeds

![request](gifs/request.gif)

## Scan

    networker scan --host farishuskovic.dev

## Capture

    networker capture -d en0 -s 10

![cap](gifs/cap.gif)

![scan](gifs/scan.gif)