## networker scan

Scan hosts for open ports.

```
networker scan [flags]
```

### Examples

```

Scan well-known ports of single device on network:
networker scan 127.0.0.1

Scan well-known ports of all devices on network:
networker scan

Scan all ports of single device on network:
networker scan 127.0.0.1 --all-ports

Output a scan as json:
networker scan 127.0.0.1 --json

```

### Options

```
      --all-ports   Scan all ports(scans first 1024 if not enabled).
  -h, --help        help for scan
      --json        Output as json.
```

### SEE ALSO

* [networker](networker.md)	 - A simple networking utility.

###### Auto generated by spf13/cobra on 2-Apr-2022