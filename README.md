# Networker

A practical CLI tool for network administration.

# Installation

From project root...

    go install .

# Exploring Commands

- [List](#list)
- [Lookup](#lookup)
- [Capture](#capture)
- [Scan](#scan)


# General Usage

```
Usage:
    networker [flags]
    networker [command]

Available Commands:
    capture     capture network packets on specified devices.
    list        list information on connected device(s).
    lookup      lookup hostnames, IP's, MX records, nameservers, and general network information.
    scan        scan for exposed ports on a designated IP.
    help        Help about any command

Flags:
    -h, --help   help for networker

Use "networker [command] --help" for more information about a command.

```

# List

    list information on connected device(s).

    Usage:
    networker list [flags]

    Aliases:
    list, ls

    Examples:

    networker ls --me -a

    Flags:
    -a, --all    enable this to list all connected network interface devices and associated information
    -h, --help   help for list
        --me     enable this to list the name, local IP, remote IP, and router IP for this machine



# Lookup

    lookup hostnames, IP addresses, MX records, nameservers, and general network information.

    Usage:
    networker lookup [flags]

    Aliases:
    lookup, lu

    Examples:

    lookup network : networker lookup --network www.farishuskovic.dev
    lookup hostname : networker lookup --hostnames 192.81.212.192
    lookup nameserver : networker lookup --nameservers farishuskovic.dev
    lookup ip : networker lookup --addresses farishuskovic.dev


    Flags:
    -a, --addresses string     look up IP addresses by hostname
    -h, --help                 help for lookup
        --hostnames string     look up hostnames by IP address
    -m, --mx string            look up MX records by domain
    -n, --nameservers string   look up name server by hostname
        --network string       look up the network for a hostname

# Capture

    capture network packets on specified devices.

    Usage:
    networker capture [flags]

    Aliases:
    capture, c, cap

    Examples:
    capture pkts on en1 for 10s or until 100 pkts captured:
    networker capture --devices en1 --seconds 10 --out myCaptureSession --limit --num 100 --verbose

    short form: networker c -d en1 -s 10 -o myCaptureSession -l -n 100 -v


    Flags:
    -d, --devices strings   devices on which to capture network packets (comma separated).
    -h, --help              help for capture
    -l, --limit             enable packet capture limiting(must use with --num || -n to specify number).
    -n, --num int           number of packets to capture (accumulative for all devices)
    -o, --out string        specify outfile to write captured packets to
    -s, --seconds int       Amount of seconds to run capture
    -v, --verbose           enable verbose logging.


# Scan

    scan for exposed ports on a designated IP.

    Usage:
    networker scan [flags]

    Aliases:
    scan, s

    Examples:
    scan only a specified set of TCP ports and only log if they're open:
    networker scan --ip <someIPaddress> --ports 22,80,3389 --open-only

    scan all TCP ports up to port 1024 and only log status if they're open:
    networker scan --ip <someIPaddress> --up-to 1024 --tcp-only --open-only


    short form: networker s --ip <someIPaddress> --up-to 1024 -t -o



    Flags:
    -h, --help         help for scan      --ip string    IP address to scan
    -o, --open-only    enable to only log open ports
    -p, --ports ints   explicitly specify which ports you want scanned (comma separated). If not specifi
    ed, all ports will be scanned unless up-to flag is specified.
    -t, --tcp-only     enable to scan only tcp ports
        --udp-only     enable to scan only udp ports
    -u, --up-to int    scan all ports up to a specified value