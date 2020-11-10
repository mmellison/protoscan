# ProtoScan

A primitive, proof of concept protocol scanner.
ProtoScan can be used to scrape protocol and version information from different servers and report 
results in a simple-to-read JSON format.

## Quickstart

### Build

The included makefile can be used to quickly get up and running.

```
> make
```

This will produce the scanner binary `bin/scanner`. If you cannot use the included file, you can build directly:

```
> CGO_ENABLED=0 go build -o bin/scanner ./cmd/scanner
```

### Usage

```
> scanner localhost:3306
```

See the help for additional options.

```
> scanner --help
usage: scanner [<flags>] [<target>]

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man).
  --init-timeout=10s  Maximum amount of time to wait for a connection to be made
  --read-timeout=5s   Maximum amount of time to wait for server to respond once a connection is made. Set to 0 to wait indefinitely.

Args:
  [<target>]  Target host and port to scan

```

## Supported Protocols

ProtoScan currently supports the following protocols and versions:

### MySQL/MariaDB

Supports reporting information from a MySQL Protocol handshake.
Supported protocol versions:
- v10
- v9

Example report:

```json
{
  "target": "127.0.0.1:3306",
  "when": "2020-11-10T15:14:39.011175257-05:00",
  "proto_version": 10,
  "handshake": {
    "server_version": "5.5.5-10.0.30-MariaDB",
    "thread_id": 164578,
    "auth_plugin_data": "aGthJCxPLS9GUUVGb3tAaXIiTT4A",
    "character_set": 8,
    "capability_flags": [
      "CLIENT_LONG_PASSWORD",
      "CLIENT_FOUND_ROWS",
      "CLIENT_LONG_FLAG",
      "CLIENT_CONNECT_WITH_DB",
      "CLIENT_NO_SCHEMA",
      "CLIENT_COMPRESS",
      "CLIENT_TRANSACTIONS",
      "CLIENT_RESERVED2",
      "CLIENT_MULTI_STATEMENTS",
      "CLIENT_MULTI_RESULTS",
      "CLIENT_PS_MULTI_RESULTS",
      "CLIENT_PLUGIN_AUTH",
      "CLIENT_CONNECT_ATTRS",
      "CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA",
      "CLIENT_CAN_HANDLED_EXPIRED_PASSWORDS",
      "CLIENT_SESSION_TRACK",
      "CLIENT_DEPRECATE_EOF",
      "CLIENT_OPTIONAL_RESULTSET_METADATA",
      "CLIENT_SSL_VERIFY_SERVER_CERT",
      "CLIENT_REMEMBER_OPTIONS"
    ],
    "server_status_flags": [
      "SERVER_STATUS_AUTOCOMMIT"
    ],
    "auth_plugin_name": "mysql_native_password"
  }
}
```