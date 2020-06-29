# Lookup

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
