# traefik-certs-dumper

[![GitHub release](https://img.shields.io/github/release/ldez/traefik-certs-dumper.svg)](https://github.com/ldez/traefik-certs-dumper/releases/latest)
[![Build Status](https://travis-ci.org/ldez/traefik-certs-dumper.svg?branch=master)](https://travis-ci.org/ldez/traefik-certs-dumper)
[![Docker Information](https://images.microbadger.com/badges/image/ldez/traefik-certs-dumper.svg)](https://hub.docker.com/r/ldez/traefik-certs-dumper/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/traefik-certs-dumper)](https://goreportcard.com/report/github.com/ldez/traefik-certs-dumper)

```yaml
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

```yaml
Dump the content of the "acme.json" file from Traefik to certificates.

Usage:
  traefik-certs-dumper dump [flags]

Flags:
      --crt-ext string    The file extension of the generated certificates. (default ".crt")
      --crt-name string   The file name (without extension) of the generated certificates. (default "certificate")
      --dest string       Path to store the dump content. (default "./dump")
      --domain-subdir     Use domain as sub-directory.
  -h, --help              help for dump
      --key-ext string    The file extension of the generated private keys. (default ".key")
      --key-name string   The file name (without extension) of the generated private keys. (default "privatekey")
      --source string     Path to 'acme.json' file. (default "./acme.json")
```

## Examples

### Simple Dump

```console
$ traefik-certs-dumper dump
dump
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key

```

#### Change source and destination

```console
$ traefik-certs-dumper dump --source ./acme.json --dest ./dump/test
test
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key

```

### Use domain as sub-directory

```console
$ traefik-certs-dumper dump --domain-subdir=true
dump
├──my.domain.com
│  ├──certificate.crt
│  └──privatekey.key
└──private
   └──letsencrypt.key
```

#### Change file extension

```console
$ traefik-certs-dumper dump --domain-subdir=true --crt-ext=.pem --key-ext=.pem
dump
├──my.domain.com
│  ├──certificate.pem
│  └──privatekey.pem
└──private
   └──letsencrypt.key
```

#### Change file name

```console
$ traefik-certs-dumper dump --domain-subdir=true --crt-name=fullchain --key-name=privkey
dump
├──my.domain.com
│  ├──fullchain.crt
│  └──privkey.key
└──private
   └──letsencrypt.key
```
