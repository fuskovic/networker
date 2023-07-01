## networker shell serve

Start a shell server.

```
networker shell serve [flags]
```

### Examples

```

# Serve using the defaults(bash is the default shell and 4444 is the default port):

	nw shell serve

# Serve a particular shell on a particular port:

	nw shell serve zsh -p 9000


```

### Options

```
  -h, --help       help for serve
  -p, --port int   Port to serve shell on. (default 4444)
```

### Options inherited from parent commands

```
  -o, --output string   Output format. Supported values include json and yaml.
```

### SEE ALSO

* [networker shell](networker_shell.md)	 - Serve and establish connections with remote shells.

###### Auto generated by spf13/cobra on 4-Jun-2022