
# Scan

`Warning` : This scanner is noisy.

## Usage

    Scan a host for exposed ports.

    Flags:
    -h, --help         help for scan
        --ip string    IP address to scan.
    -o, --open-only    Only print the ports that are open.
    -p, --ports ints   Specify a comma-separated list of ports to scan. (scans all ports if left unspecified)
    -t, --tcp-only     Only scan TCP ports.
        --udp-only     Only scan UDP ports.
    -u, --up-to int    Scan all ports up to a given port number.


## Examples:

Scan three explicit TCP ports.

    networker scan --ip <host> --ports 22,80,3389 --tcp-only

Scan all TCP ports but only print the ones that are open.

    networker scan --ip <host> --tcp-only --open-only

Scan all ports up to a certain port number but only print the TCP ports that are open.

    networker scan--ip <host> --up-to 1024 --open-only --tcp-only


