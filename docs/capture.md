# Capture

## Index

|Name|Description|Examples|
|---|---|---|
|request|Send an HTTP request.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/request.md)|
|list|List information on connected network devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/list.md)|
|lookup|Lookup hostnames, IP addresses, nameservers, and networks.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/lookup.md)|
|proxy|Proxy ingress to an upstream server.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/proxy.md)|
|capture|Capture network packets on specified devices.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/capture.md)|
|scan|Scan a host for exposed ports.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/scan.md)|
|backdoor (Unsafe)|Serve shell access over TCP and connect remotely.|[Docs](https://github.com/fuskovic/networker/tree/master/docs/backdoor.md)|

## Usage

    Capture network packets on specified devices.

    Flags:
    -d, --devices strings   Comma-separated list of devices to capture packets on.
    -h, --help              help for capture
    -l, --limit             Limit the number of packets to capture. (must be used with the --num flag)
    -n, --num int           Number of total packets to capture across all devices.
    -o, --out string        Name of an output file to write the packets to.
    -s, --seconds int       Amount of seconds to run capture for.
    -v, --verbose           Enable verbose logging.

## Examples


Write a `10s` capture of `en0` to stdout.


    networker capture --devices en0 --seconds 10



Write the capture to an outfile.


    networker capture --devices en0 --seconds 10 --out myCaptureSession
