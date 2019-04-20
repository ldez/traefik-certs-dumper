# Integration testing

## Preperation
1.  Create valid ACME file `/tmp/acme.json`
1.  Start backends using docker
    ```console
    docker run -d -p 8500:8500 --name consul consul
    docker run -d -p 2181:2181 --name zookeeper zookeeper
    docker run -d -p 2379:2379 --name etcd quay.io/coreos/etcd:v3.3.12 etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2380
    ```
1.  Build tests
    ```console
    export GO111MODULE=on
    go build
    ```
1.  Initialize backends
    ```console
    ./tests
    ```
1.  Run certs dumper without watching
    ```console
    ../traefik-certs-dumper dump --source.file=/tmp/acme.json
    ../traefik-certs-dumper dump --source consul --source.kv.endpoints=localhost:8500
    ../traefik-certs-dumper dump --source etcd --source.kv.endpoints=localhost:2379
    ../traefik-certs-dumper dump --source boltdb --source.kv.endpoints=/tmp/my.db
    ../traefik-certs-dumper dump --source zookeeper --source.kv.endpoints=localhost:2181
    ```
1.  Run certs dumper with watching
    ```console
    ../traefik-certs-dumper dump --watch --source.file=/tmp/acme.json
    ../traefik-certs-dumper dump --watch --source consul --source.kv.endpoints=localhost:8500
    ../traefik-certs-dumper dump --watch --source etcd --source.kv.endpoints=localhost:2379
    ../traefik-certs-dumper dump --watch --source zookeeper --source.kv.endpoints=localhost:2181
    ```
    While watching is enabled, run `./tests` again for KV backends or manipulate `/tmp/acme.json` for file backend that change events are triggered.
