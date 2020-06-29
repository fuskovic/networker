# Backdoor

`Warning` : This command is unsafe right now because the shell session is not safely being terminated. Don't use this command for now.

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

    Serve shell access over TCP and connect remotely.

    Flags:
    -a, --address string   Address to connect to. (format: <host>:<port>)
        --connect          Enable connect mode. (must be used with the --address flag)
        --create           Enable create mode. (must be used with the --port flag)
    -h, --help             help for backdoor
    -p, --port int         Port number to serve shell access on. (format: 80)

## Examples

Serve shell access on server A.

    networker backdoor --create --port <port>

From client A, use networker to gain shell access on server A.

    networker backdoor --connect --address <host>:<port>