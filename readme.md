# traefik-certs-dumper

[![Build Status](https://travis-ci.org/ldez/traefik-certs-dumper.svg?branch=master)](https://travis-ci.org/ldez/traefik-certs-dumper)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/traefik-certs-dumper)](https://goreportcard.com/report/github.com/ldez/traefik-certs-dumper)


```
Dump the content of the "acme.json" file from Traefik to certificates.

Usage:
  traefik-certs-dumper [command]

Available Commands:
  dump        Dump Let's Encrypt certificates from Traefik
  help        Help about any command
  version     Display version

Flags:
  -h, --help      help for certs-dumper
      --version   version for certs-dumper

Use "traefik-certs-dumper [command] --help" for more information about a command.
```

```
Dump the content of the "acme.json" file from Traefik to certificates.

Usage:
  traefik-certs-dumper dump [flags]

Flags:
      --crt-ext string   The file extension of the generated certificates. (default ".crt")
      --dest string      Path to store the dump content. (default "./dump")
  -h, --help             help for dump
      --key-ext string   The file extension of the generated private keys. (default ".key")
      --source string    Path to 'acme.json' file. (default "./acme.json")
      --use-subdir       Use separated directories for certificates and keys. (default true)
```

## Examples

```bash
traefik-certs-dumper dump
```

```bash
traefik-certs-dumper dump --source ./acme.json --dest ./dump
```

```bash
traefik-certs-dumper dump --crt-ext=.pem --key-ext=.pem
```

```bash
traefik-certs-dumper dump --use-subdir=false
```

- https://github.com/containous/traefik/issues/4381
- https://github.com/containous/traefik/issues/2418
- https://github.com/containous/traefik/issues/3847
- https://github.com/SvenDowideit/traefik-certdumper
