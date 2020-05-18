# traefik-certs-dumper

[![GitHub release](https://img.shields.io/github/release/ldez/traefik-certs-dumper.svg)](https://github.com/ldez/traefik-certs-dumper/releases/latest)
[![Build Status](https://travis-ci.com/ldez/traefik-certs-dumper.svg?branch=master)](https://travis-ci.com/ldez/traefik-certs-dumper)
[![Docker Information](https://images.microbadger.com/badges/image/ldez/traefik-certs-dumper.svg)](https://hub.docker.com/r/ldez/traefik-certs-dumper/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/traefik-certs-dumper)](https://goreportcard.com/report/github.com/ldez/traefik-certs-dumper)

If you appreciate this project:

[![Sponsor](https://img.shields.io/badge/Sponsor%20me-%E2%9D%A4%EF%B8%8F-pink)](https://github.com/sponsors/ldez)

## Features

- Supported sources:
    - file ("acme.json")
    - KV stores (Consul, Etcd, Zookeeper, Boltdb)
- Watch changes:
    - from file ("acme.json")
    - from KV stores (Consul, Etcd, Zookeeper)
- Output formats:
    - use domain as sub-directory (allow custom names and extensions)
    - flat (domain as filename)
- Hook (only with watch mode and if the data source changes)

## Installation

### Download / CI Integration

```bash
curl -sfL https://raw.githubusercontent.com/ldez/traefik-certs-dumper/master/godownloader.sh | bash -s -- -b $GOPATH/bin v2.7.0
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

- [traefik-certs-dumper](docs/traefik-certs-dumper.md)
- [traefik-certs-dumper file](docs/traefik-certs-dumper_file.md)
- [traefik-certs-dumper kv](docs/traefik-certs-dumper_kv.md)

## Examples

**Note:** to dump data from Traefik v2, the CLI flag `--version v2` must be added.

### `docker-compose`

`docker-compose.yml`:

```yaml
version: '3'

services:
    traefik:
        image: traefik:1.7
        command:
            --entryPoints='Name:http Address::80'
            --entryPoints='Name:https Address::443 TLS'
            --defaultentrypoints=http,https
            --logLevel=DEBUG
            --docker
            --docker.exposedByDefault=false
            --acme
            --acme.acmeLogging=true
            --acme.entrypoint=https
            --acme.storage=/data/acme.json
            --acme.onHostRule=true
            --acme.httpChallenge.entryPoint=http
        ports:
            - 8001:80
            - 8002:443
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
            - .:/data

    traefik-certs-dumper:
        image: ldez/traefik-certs-dumper:v2.7.0
        entrypoint: sh -c '
            apk add jq
            ; while ! [ -e /data/acme.json ]
                || ! [ `jq ".Certificates | length" /data/acme.json` != 0 ]; do
                    sleep 1
                ; done
            && traefik-certs-dumper file --watch 
                --source /data/acme.json --dest /data/certs'
        volumes:
            - .:/data

    whoami:
        image: containous/whoami
        labels:
            traefik.enable: true
            traefik.frontend.rule: Host:example.com
```

`traefik` ports are published to `8001` and `8002`. It's assumed here that you need certificates as separate files because you want to put `traefik` behind another proxy.

### Simple Dump

```console
$ traefik-certs-dumper file
dump
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key
```

### Change source and destination

```console
$ traefik-certs-dumper file --source ./acme.json --dest ./dump/test
test
├──certs
│  └──my.domain.com.key
└──private
   ├──my.domain.com.crt
   └──letsencrypt.key
```

### Use domain as sub-directory

```console
$ traefik-certs-dumper file --domain-subdir=true
dump
├──my.domain.com
│  ├──certificate.crt
│  └──privatekey.key
└──private
   └──letsencrypt.key
```

#### Change file extension

```console
$ traefik-certs-dumper file --domain-subdir --crt-ext=.pem --key-ext=.pem
dump
├──my.domain.com
│  ├──certificate.pem
│  └──privatekey.pem
└──private
   └──letsencrypt.key
```

#### Change file name

```console
$ traefik-certs-dumper file --domain-subdir --crt-name=fullchain --key-name=privkey
dump
├──my.domain.com
│  ├──fullchain.crt
│  └──privkey.key
└──private
   └──letsencrypt.key
```

### KV store

#### Consul

```console
$ traefik-certs-dumper kv consul --endpoints localhost:8500
```

#### Etcd

```console
$ traefik-certs-dumper kv etcd --endpoints localhost:2379
```

#### Boltdb

```console
$ traefik-certs-dumper kv boltdb --endpoints /the/path/to/mydb.db
```

#### Zookeeper

```console
$ traefik-certs-dumper kv zookeeper --endpoints localhost:2181
```


