# Lookup

## Index

|Name|Description|Examples|Code|
|---|---|---|---|
|request|Send an HTTP request.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/request.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/request.go)|
|list|List information on connected network devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/list.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/list.go)|
|lookup|Lookup hostnames, IP addresses, nameservers, and networks.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/lookup.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/lookup.go)|
|proxy|Proxy ingress to an upstream server.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/proxy.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/proxy.go)|
|capture|Capture network packets on specified devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/capture.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/capture.go)|
|scan|Scan a host for exposed ports.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/scan.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/scan.go)|
|backdoor (Unsafe)|Serve shell access over TCP and connect remotely.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/backdoor.md)|[File](https://github.com/fuskovic/networker/tree/master/cmd/backdoor.go)|

## Usage

    Lookup hostnames, IP addresses, nameservers, and networks.

    Flags:
    -a, --addresses string     Look up IP addresses for a given hostname.
    -h, --help                 help for lookup
        --hostnames string     Look up hostnames for a given IP address.
    -s, --nameservers string   Look up nameservers for a given hostname.
    -n, --network string       Look up the network a given hostname belongs to.


## Examples

Look up the network for a given host.

    networker lookup --network 31.13.65.36

Look up hostnames for a given IP.

    networker lookup --hostnames 157.240.195.35

Look up nameservers for a given hostname.

    networker lookup --nameservers youtube.com

Look up addresses for a given hostname.

    networker lookup --addresses youtube.com
