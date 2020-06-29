# List

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
