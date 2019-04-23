# Integration testing

## Preparation

- Create valid ACME file `./acme.json`

- Start backends using docker

```bash
docker-compose -f integrationtest/docker-compose.yml up
```

- Initialize backends

```bash
go run integrationtest/loader.go
```

## Run certs dumper without watching

```bash
# http://localhost:8500/ui/
traefik-certs-dumper kv consul --endpoints localhost:8500

traefik-certs-dumper kv etcd --endpoints localhost:2379

traefik-certs-dumper kv boltdb --endpoints /tmp/test-traefik-certs-dumper.db

traefik-certs-dumper kv zookeeper --endpoints localhost:2181
```

## Run certs dumper with watching

While watching is enabled, run `loader.go` again for KV backends or manipulate `./acme.json` for file backend that change events are triggered.

```bash
traefik-certs-dumper kv consul --watch --endpoints localhost:8500

traefik-certs-dumper kv etcd --watch --endpoints localhost:2379

traefik-certs-dumper kv zookeeper --watch --endpoints localhost:2181
```
