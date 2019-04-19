# traefik-certs-dumper

[![GitHub release](https://img.shields.io/github/release/ldez/traefik-certs-dumper.svg)](https://github.com/ldez/traefik-certs-dumper/releases/latest)
[![Build Status](https://travis-ci.org/ldez/traefik-certs-dumper.svg?branch=master)](https://travis-ci.org/ldez/traefik-certs-dumper)
[![Docker Information](https://images.microbadger.com/badges/image/ldez/traefik-certs-dumper.svg)](https://hub.docker.com/r/ldez/traefik-certs-dumper/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/traefik-certs-dumper)](https://goreportcard.com/report/github.com/ldez/traefik-certs-dumper)

[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/ldez)

## Installation

### Download / CI Integration

```bash
curl -sfL https://raw.githubusercontent.com/ldez/traefik-certs-dumper/master/godownloader.sh | bash -s -- -b $GOPATH/bin v1.5.0
```

<!--
To generate the script:

```bash
godownloader --repo=ldez/traefik-certs-dumper -o godownloader.sh

# or

godownloader --repo=ldez/traefik-certs-dumper > godownloader.sh
```
-->

### From Binaries

You can use pre-compiled binaries:

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/traefik-certs-dumper/releases/)
* Unzip the archive.
* Add `traefik-certs-dumper` in your `PATH`.

### From Docker

```bash
docker run ldez/traefik-certs-dumper:<tag_name>
```

## Usage

```yaml
Dump ACME data from Traefik of different storage backends to certificates.

Usage:
  traefik-certs-dumper [command]

Available Commands:
  dump        Dump Let's Encrypt certificates from Traefik
  help        Help about any command
  version     Display version

Flags:
  -h, --help      help for traefik-certs-dumper
      --version   version for traefik-certs-dumper

Use "traefik-certs-dumper [command] --help" for more information about a command.
```

```yaml
Dump ACME data from Traefik of different storage backends to certificates.

Usage:
  traefik-certs-dumper dump [flags]

Flags:
      --crt-ext string                        The file extension of the generated certificates. (default ".crt")
      --crt-name string                       The file name (without extension) of the generated certificates. (default "certificate")
      --dest string                           Path to store the dump content. (default "./dump")
      --domain-subdir                         Use domain as sub-directory.
  -h, --help                                  help for dump
      --key-ext string                        The file extension of the generated private keys. (default ".key")
      --key-name string                       The file name (without extension) of the generated private keys. (default "privatekey")
      --source source.<type>.                 Source type, one of 'file', 'consul', 'etcd', 'zookeeper', 'boltdb'. Options for each source type are prefixed with source.<type>. (default "file")
      --source.file string                    Path to 'acme.json' for file source. (default "./acme.json")
      --source.kv.boltdb.bucket string        Bucket for boltdb. (default "traefik")
      --source.kv.boltdb.persist-connection   Persist connection for boltdb.
      --source.kv.connection-timeout int      Connection timeout in seconds.
      --source.kv.consul.token string         Token for consul.
      --source.kv.endpoints string            Comma seperated list of endpoints. (default "localhost:8500")
      --source.kv.etcd.sync-period int        Sync period for etcd in seconds.
      --source.kv.password string             Password for connection.
      --source.kv.tls.ca-cert-file string     Root CA file for certificate verification if TLS is enabled.
      --source.kv.tls.enable                  Enable TLS encryption.
      --source.kv.tls.insecureskipverify      Trust unverified certificates if TLS is enabled.
      --source.kv.username string             Username for connection.
      --watch                                 Enable watching changes.
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

### Enabled watching

```console
$ traefik-certs-dumper dump --watch
2019/04/19 16:56:34 wrote new configuration
dump
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key
2019/04/19 16:57:14 wrote new configuration
dump
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key

```

### Consul backend

```console
$ traefik-certs-dumper dump --source consul --source.kv.endpoints=localhost:8500
```

### Etcd backend

```console
$ traefik-certs-dumper dump --source etcd --source.kv.endpoints=localhost:2379
```

### Boltdb backend

```console
$ traefik-certs-dumper dump --source boltdb --source.kv.endpoints=/tmp/my.db
```

### Zookeeper backend

```console
$ traefik-certs-dumper dump --source zookeeper --source.kv.endpoints=localhost:2181
```

### Change source and destination

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
