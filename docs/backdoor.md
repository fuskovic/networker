# Backdoor

`Warning` : This command is unsafe right now because the shell session is not safely being terminated. Don't use this command for now.

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