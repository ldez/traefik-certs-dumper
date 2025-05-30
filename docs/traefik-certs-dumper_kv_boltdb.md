## traefik-certs-dumper kv boltdb

Dump the content of BoltDB.

### Synopsis

Dump the content of BoltDB.

```
traefik-certs-dumper kv boltdb [flags]
```

### Options

```
      --bucket string        Bucket for boltdb. (default "traefik")
  -h, --help                 help for boltdb
      --persist-connection   Persist connection for boltdb.
```

### Options inherited from parent commands

```
      --clean                    Clean destination folder before dumping content. (default true)
      --config string            config file (default is $HOME/.traefik-certs-dumper.yaml)
      --connection-timeout int   Connection timeout in seconds.
      --crt-ext string           The file extension of the generated certificates. (default ".crt")
      --crt-name string          The file name (without extension) of the generated certificates. (default "certificate")
      --dest string              Path to store the dump content. (default "./dump")
      --domain-subdir            Use domain as sub-directory.
      --endpoints strings        List of endpoints. (default [localhost:8500])
      --key-ext string           The file extension of the generated private keys. (default ".key")
      --key-name string          The file name (without extension) of the generated private keys. (default "privatekey")
      --password string          Password for connection.
      --post-hook string         Execute a command only if changes occurs on the data source. (works only with the watch mode)
      --prefix string            Prefix used for KV store. (default "traefik")
      --suffix string            Suffix/Storage used for KV store. (default "/acme/account/object")
      --tls                      Enable TLS encryption.
      --tls.ca string            Root CA for certificate verification if TLS is enabled
      --tls.ca.optional          
      --tls.cert string          TLS cert
      --tls.insecureskipverify   Trust unverified certificates if TLS is enabled.
      --tls.key string           TLS key
      --username string          Username for connection.
      --watch                    Enable watching changes.
```

### SEE ALSO

* [traefik-certs-dumper kv](traefik-certs-dumper_kv.md)	 - Dump the content of a KV store.

###### Auto generated by spf13/cobra on 21-Feb-2025
