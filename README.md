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
- [Proxy](#proxy)
- [Backdoor](#backdoor)


# General Usage

```
Usage:
    networker [flags]
    networker [command]

Available Commands:
    backdoor    create and connect to backdoors to gain shell access over TCP
    capture     capture network packets on specified devices.
    help        Help about any command
    list        list information on connected device(s).
    lookup      lookup hostnames, IP addresses, nameservers, and general network information.
    proxy       forward network traffic from one network connection to another
    scan        scan for exposed ports on a designated IP.

Flags:
    -h, --help   help for networker

Use "networker [command] --help" for more information about a command.

```

# List

    list information on connected device(s).
    Must be run as root.

    Usage:
        networker list [flags]

    Aliases:
        list, ls

    Examples:

        sudo networker ls --me -a

    Flags:
        -a, --all    enable this to list all connected network interface devices and associated information(must be run as root)
        -h, --help   help for list
            --me     enable this to list the name, local IP, remote IP, and router IP for this machine

# Lookup

    lookup hostnames, IP addresses, nameservers, and general network information.

    Usage:
        networker lookup [flags]

    Aliases:
        lookup, lu

    Examples:

    lookup network : 
        networker lookup --network facebook.com || 31.13.65.36

    lookup hostname : 
        networker lookup --hostnames 157.240.195.35

    lookup nameserver : 
        networker lookup --nameservers youtube.com

    lookup ip : 
        networker lookup --addresses youtube.com



    Flags:
        -a, --addresses string     look up IP addresses by hostname
        -h, --help                 help for lookup
            --hostnames string     look up hostnames by IP address
        -n, --nameservers string   look up name servers by hostname
            --network string       look up the network of a host


# Capture

    capture network packets on specified devices.

    Usage:
        networker capture [flags]

    Aliases:
        capture, c, cap

    Examples:
        capture pkts on en1 for 10s or until 100 pkts captured:
            networker capture --devices en1 --seconds 10 --out myCaptureSession --limit --num 100 --verbose

        short form: 
            networker c -d en1 -s 10 -o myCaptureSession -l -n 100 -v


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

        short form: 
            networker s --ip <someIPaddress> --up-to 1024 -t -o


    Flags:
        -h, --help         help for scan      --ip string    IP address to scan
        -o, --open-only    enable to only log open ports
        -p, --ports ints   explicitly specify which ports you want scanned (comma separated). If not specifi
        ed, all ports will be scanned unless up-to flag is specified.
        -t, --tcp-only     enable to scan only tcp ports
            --udp-only     enable to scan only udp ports
        -u, --up-to int    scan all ports up to a specified value


# Proxy

    forward network traffic from one network connection to another

    Usage:
        networker proxy [flags]

    Aliases:
        proxy, p

    Examples:

        long format:

            networker proxy --listen-on <port> -upstream <host>:<port>

        short format:

            networker p -l <port> -u <host>:<port>

    Flags:
        -h, --help              help for proxy
        -l, --listen-on int     port for proxy to listen on
        -u, --upstream string   <host>:<port> to proxy traffic to

# Backdoor

    create and connect to backdoors to gain shell access over TCP

    Usage:
        networker backdoor [flags]

    Aliases:
        backdoor, bd, b

    Examples:

        long format:

            networker backdoor --create --port 4444

        short format:

            networker backdoor --connect --address <host>:4444

    Flags:
        -a, --address string   address of remote target to connect to(format: <host>:<port>))
            --connect          connect to a TCP backdoor(must be used with --address flag)
            --create           create a TCP backdoor(must be used with --port flag)
        -h, --help             help for backdoor
        -p, --port int         port number to listen for connections on