# Networker

[![Go Report Card](https://goreportcard.com/badge/github.com/fuskovic/networker)](https://goreportcard.com/report/github.com/fuskovic/networker)

A practical CLI tool for network administration.

# Installation

    go get -u github.com/fuskovic/networker

# Explore Commands

|Name|Description|Examples|Code|
|---|---|---|---|
|request|Send an HTTP request.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/request.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/request.go)|
|list|List information on connected network devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/list.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/list.go)|
|lookup|Lookup hostnames, IP addresses, nameservers, and networks.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/lookup.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/lookup.go)|
|proxy|Proxy ingress to an upstream server.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/proxy.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/proxy.go)|
|capture|Capture network packets on specified devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/capture.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/capture.go)|
|scan|Scan a host for exposed ports.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/scan.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/scan.go)|
|backdoor (Unsafe)|Serve shell access over TCP and connect remotely.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/backdoor.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/backdoor.go)|