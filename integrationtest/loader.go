package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kvtools/boltdb"
	"github.com/kvtools/consul"
	"github.com/kvtools/etcdv3"
	"github.com/kvtools/valkeyrie"
	"github.com/kvtools/zookeeper"
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

	err = putData(ctx, consul.StoreName, []string{"localhost:8500"},
		&consul.Config{ConnectionTimeout: 3 * time.Second}, content)
	if err != nil {
		return err
	}

	// ETCD v3
	err = putData(ctx, etcdv3.StoreName, []string{"localhost:2379"},
		&etcdv3.Config{ConnectionTimeout: 3 * time.Second}, content)
	if err != nil {
		return err
	}

	// Zookeeper
	err = putData(ctx, zookeeper.StoreName, []string{"localhost:2181"},
		&zookeeper.Config{ConnectionTimeout: 3 * time.Second}, content)
	if err != nil {
		return err
	}

	// BoltDB
	err = putData(ctx, boltdb.StoreName, []string{"/tmp/test-traefik-certs-dumper.db"},
		&boltdb.Config{ConnectionTimeout: 3 * time.Second, Bucket: "traefik"}, content)
	if err != nil {
		return err
	}

	return nil
}

func putData(ctx context.Context, backend string, addrs []string, options valkeyrie.Config, content []byte) error {
	kvStore, err := valkeyrie.NewStore(ctx, backend, addrs, options)
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
