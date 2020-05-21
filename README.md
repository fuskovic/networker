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
A practical CLI tool for network administration.

Usage:
  networker [flags]
  networker [command]

Available Commands:
  backdoor    Create and connect to backdoors to gain shell access over TCP.
  capture     Capture network packets on specified devices.
  help        Help about any command
  list        List information on connected network devices.
  lookup      Lookup hostnames, IP addresses, nameservers, and general network information.
  proxy       Forward network traffic from one network connection to another.
  scan        Scan an IP for exposed ports.

Flags:
  -h, --help   help for networker

Use "networker [command] --help" for more information about a command.

```

# List

    List information on connected network devices.

    Usage:
    networker list [flags]

    Aliases:
    list, ls

    Examples:

    List the IP of the current network gateway, local IP of this machine, and remote IP of this machine:

            long form:

                    networker list --me

            short form:

                    networker ls -m

    List the hostname, IP address, and connection status of all devices on the current network:

            long form:

                    sudo networker list --all

            short form:

                    sudo networker ls -a


    Flags:
    -a, --all    List the IP, hostname, and connection status of all devices on this network. (must be run as root)
    -h, --help   help for list
    -m, --me     List the name, local IP, remote IP, and router IP for this machine and the network it's connected to.

# Lookup

    Lookup hostnames, IP addresses, nameservers, and general network information.

    Usage:
    networker lookup [flags]

    Aliases:
    lookup, lu

    Examples:

    Look up the network for a given hostname or IP:

            long form:

                    networker lookup --network 31.13.65.36

            short form:

                    networker lu -n 31.13.65.36

    Look up the hostname for a given IP:

            long form:

                    networker lookup --hostnames 157.240.195.35

            short form:

                    no short form as -h is reserved for help

    Look up nameservers for a given hostname:

            long form:

                    networker lookup --nameservers youtube.com

            short form:

                    networker lu -s youtube.com

    Look up the addresses for a given hostname:

            long form:

                    networker lookup --addresses youtube.com

            short form:

                    networker lu -a youtube.com


    Flags:
    -a, --addresses string     Look up IP addresses for a given hostname.
    -h, --help                 help for lookup
        --hostnames string     Look up hostnames for a given IP address.
    -s, --nameservers string   Look up nameservers for a given hostname.
    -n, --network string       Look up the network a given hostname belongs to.


# Capture

    Capture network packets on specified devices.

    Usage:
    networker capture [flags]

    Aliases:
    capture, c, cap

    Examples:

    Capture packets on en1 for 10 seconds or until 100 packets have been captured and log the capture status during capture:

            long form:

                    networker capture --devices en1 --seconds 10 --out myCaptureSession --limit --num 100 --verbose

            short form:

                    networker c -d en1 -s 10 -out myCaptureSession -l -n 100 -v


    Flags:
    -d, --devices strings   Comma-separated list of devices to capture packets on.
    -h, --help              help for capture
    -l, --limit             Limit the number of packets to capture. (must be used with the --num flag)
    -n, --num int           Number of total packets to capture across all devices.
    -o, --out string        Name of an output file to write the packets to.
    -s, --seconds int       Amount of seconds to run capture for.
    -v, --verbose           Enable verbose logging.


# Scan

    Scan an IP for exposed ports.

    Usage:
    networker scan [flags]

    Aliases:
    scan, s

    Examples:

    Scan a comma-separated list of TCP ports of an address and only log out the ones that are open:

            long form:

                    networker scan --ip <address> --ports 22,80,3389 --tcp-only --open-only

            short form:

                    networker s --ip <address> -p 22,80,3389 -t -o

    Scan all TCP ports up to port 1024 and only log out the ones that are open:

            long form:

                    networker scan --ip <someIPaddress> --up-to 1024 --tcp-only --open-only

            short form:

                    networker s --ip <address> -u 1024 -t -o


    Flags:
    -h, --help         help for scan
        --ip string    IP address to scan.
    -o, --open-only    Only print the ports that are open.
    -p, --ports ints   Specify a comma-separated list of ports to scan. (scans all ports if left unspecified)
    -t, --tcp-only     Only scan TCP ports.
        --udp-only     Only scan UDP ports.
    -u, --up-to int    Scan all ports up to a given port number.


# Proxy

    Forward network traffic from one network connection to another.

    Usage:
    networker proxy [flags]

    Aliases:
    proxy, p

    Examples:

    Start a new proxy server that listens on a given port and forwards traffic to a given address:

            long form:

                    networker proxy --listen-on <port> --upstream <host>:<port>

            short form:

                    networker p -l <port> -u <host>:<port>


    Flags:
    -h, --help              help for proxy
    -l, --listen-on int     Port to listen on.
    -u, --upstream string   Address of server to forward traffic to.

# Backdoor

    Create and connect to backdoors to gain shell access over TCP.

    Usage:
    networker backdoor [flags]

    Aliases:
    backdoor, bd, b

    Examples:

    Create a new backdoor:

            long form:

                    networker backdoor --create --port <port>

            short form:

                    networker bd --create -p <port>

    Connect to an existing backdoor:

            long form:

                    networker backdoor --connect --address <host>:<port>

            short form:

                    networker bd --connect -a <host>:<port>


    Flags:
    -a, --address string   Address of a remote target to connect to. (format: <host>:<port>)
        --connect          Connect to a TCP backdoor. (must be used with the --address flag)
        --create           Create a TCP backdoor. (must be used with the --port flag)
    -h, --help             help for backdoor
    -p, --port int         Port number to listen for connections on. (format: 80)