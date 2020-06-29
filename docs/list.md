# List

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

    List information on connected network devices.

    Flags:
    -a, --all    List the IP, hostname, and connection status of all devices on this network. (must be run as root)
    -h, --help   help for list
    -m, --me     List the local IP, remote IP, and router IP for this machine and the network it's connected to.

## Examples

A quick way to get your local and remote IP address.  Also outputs the router IP.

    networker list --me

List the hostname, IP address, and connection status of all devices on the current network. Needs to be run as root.

    sudo networker list --all
