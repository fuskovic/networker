# Networker


# Install Using Go

    go get -u github.com/fuskovic/networker/cmd/networker

# Download Pre-compiled binaries

Checkout the [releases](https://github.com/fuskovic/networker/releases) page to download the latest executables for Linux, Mac, and Windows.

# Usage 

    Usage: networker [subcommand] [flags]

    Description: A simple networking tool.

    Commands:
            ls, list         - List information on connected network devices.
            lu, lookup       - Lookup hostnames, IP addresses, nameservers, and networks.
            r, req, request  - Send an HTTP request.
            s, scan          - Scan hosts for open ports.

# Commands

## List

```
Usage: networker list [flags]

Aliases: ls

Description: List information on connected network devices.

networker list flags:
      --json   Output as json.
```


## Scan

```
Usage: networker scan [flags]

Aliases: s

Description: Scan hosts for open ports.

networker scan flags:
  -a, --all           Scan all ports(scans first 1024 if not enabled).
      --host string   Host to scan(scans all hosts on LAN if not provided).
      --json          Output as json.
```


## Lookup

    Usage: networker lookup [flags]

    Aliases: lu

    Description: Lookup hostnames, IP addresses, nameservers, and networks.

    Commands:
            hostname     - Lookup the hostname for a provided ip address.
            ip           - Lookup the ip address of the provided hostname.
            network      - Lookup the network address of a provided host.
            nameservers  - Lookup nameservers for the provided hostname.
            isp          - Lookup the internet service provider of a remote host.



## Request

    Usage: networker request [flags]

    Aliases: r, req

    Description: Send an HTTP request.

    networker request flags:
    -a, --add-headers strings   Add a list of comma-separated request headers. (format : key:value,key:value,etc...)
    -f, --file string           Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)
    -m, --method string         Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE) (default "GET")
    -t, --time-out int          Specify number of seconds for time-out. (default 3)
    -u, --url string            URL to send request.

