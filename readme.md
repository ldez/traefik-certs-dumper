# traefik-certs-dumper

[![GitHub release](https://img.shields.io/github/release/ldez/traefik-certs-dumper.svg)](https://github.com/ldez/traefik-certs-dumper/releases/latest)
[![Build Status](https://github.com/ldez/traefik-certs-dumper/workflows/Main/badge.svg?branch=master)](https://github.com/ldez/traefik-certs-dumper/actions)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/ldez/traefik-certs-dumper)](https://hub.docker.com/r/ldez/traefik-certs-dumper/)
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
curl -sfL https://raw.githubusercontent.com/ldez/traefik-certs-dumper/master/godownloader.sh | bash -s -- -b $(go env GOPATH)/bin v2.7.4
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

Examples:

- Traefik v1: [docker-compose](docs/docker-compose-traefik-v1.yml)
- Traefik v2: [docker-compose](docs/docker-compose-traefik-v2.yml)

## Usage

- [traefik-certs-dumper](docs/traefik-certs-dumper.md)
- [traefik-certs-dumper file](docs/traefik-certs-dumper_file.md)
- [traefik-certs-dumper kv](docs/traefik-certs-dumper_kv.md)

## Examples

**Note:** to dump data from Traefik v2, the CLI flag `--version v2` must be added.

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

## Hook

Hook can be a one liner passed as a string, or a file for more complex post-hook scenarios.
For the former, create a file (ex: `hook.sh`) and mount it, then pass `sh hooksh` as a parameter to `--post-hook`.

Here is a docker-compose example:

```yml
version: '3.7'

services:
# ...

  traefik-certs-dumper:
    image: ldez/traefik-certs-dumper:v2.7.4
    container_name: traefik-certs-dumper
    entrypoint: sh -c '
      apk add jq
      ; while ! [ -e /data/acme.json ]
      || ! [ `jq ".[] | .Certificates | length" /data/acme.json` != 0 ]; do
      sleep 1
      ; done
      && traefik-certs-dumper file --version v2 --watch
        --source /data/acme.json --dest /data/certs
        --post-hook "sh /hook.sh"'
    labels:
      traefik.enable: false
    volumes:
      - ./letsencrypt:/data
      - ./hook.sh:/hook.sh

# ...
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
