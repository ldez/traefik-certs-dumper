package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kvtools/valkeyrie"
	"github.com/kvtools/valkeyrie/store"
	"github.com/kvtools/valkeyrie/store/boltdb"
	"github.com/kvtools/valkeyrie/store/consul"
	etcdv3 "github.com/kvtools/valkeyrie/store/etcd/v3"
	"github.com/kvtools/valkeyrie/store/zookeeper"
)

const storeKey = "traefik/acme/account/object"

func main() {
	log.SetFlags(log.Lshortfile)

	source := "./acme.json"
	err := loadData(context.Background(), source)
	if err != nil {
		log.Fatal(err)
	}
}

func loadData(ctx context.Context, source string) error {
	content, err := readFile(source)
	if err != nil {
		return err
	}

	// Consul
	err = putData(ctx, store.CONSUL, []string{"localhost:8500"}, content)
	if err != nil {
		return err
	}

	// ETCD v3
	err = putData(ctx, store.ETCDV3, []string{"localhost:2379"}, content)
	if err != nil {
		return err
	}

	// Zookeeper
	err = putData(ctx, store.ZK, []string{"localhost:2181"}, content)
	if err != nil {
		return err
	}

	// BoltDB
	err = putData(ctx, store.BOLTDB, []string{"/tmp/test-traefik-certs-dumper.db"}, content)
	if err != nil {
		return err
	}

	return nil
}

func putData(ctx context.Context, backend store.Backend, addrs []string, content []byte) error {
	storeConfig := &store.Config{
		ConnectionTimeout: 3 * time.Second,
		Bucket:            "traefik",
	}

	switch backend {
	case store.CONSUL:
		consul.Register()
	case store.ETCDV3:
		etcdv3.Register()
	case store.ZK:
		zookeeper.Register()
	case store.BOLTDB:
		boltdb.Register()
	}

	kvStore, err := valkeyrie.NewStore(ctx, backend, addrs, storeConfig)
	if err != nil {
		return err
	}

	if err := kvStore.Put(ctx, storeKey, content, nil); err != nil {
		return err
	}

	log.Printf("Successfully updated %s.\n", backend)
	return nil
}

func readFile(source string) ([]byte, error) {
	content, err := os.ReadFile(filepath.Clean(source))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	defer func() {
		if errC := gz.Close(); errC != nil {
			log.Println(errC)
		}
	}()

	if _, err = gz.Write(content); err != nil {
		return nil, err
	}

	if err = gz.Flush(); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
