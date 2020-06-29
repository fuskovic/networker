
# Proxy

## Usage

    Proxy ingress to an upstream server.

    Flags:
    -h, --help              help for proxy
    -l, --listen-on int     Port to listen on.
    -u, --upstream string   Address of server to forward traffic to.


## Examples

Turn the current machine into a proxy server by forwarding ingress traffic on the listener to an upstream server.

    networker proxy --listen-on <port> --upstream <host>:<port>