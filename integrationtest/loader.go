package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/zookeeper"
)

const storeKey = "traefik/acme/account/object"

func main() {
	log.SetFlags(log.Lshortfile)

	source := "./acme.json"
	err := loadData(source)
	if err != nil {
		log.Fatal(err)
	}
}

func loadData(source string) error {
	content, err := readFile(source)
	if err != nil {
		return err
	}

	// Consul
	err = putData(store.CONSUL, []string{"localhost:8500"}, content)
	if err != nil {
		return err
	}

	// ETCD v3
	err = putData(store.ETCDV3, []string{"localhost:2379"}, content)
	if err != nil {
		return err
	}

	// Zookeeper
	err = putData(store.ZK, []string{"localhost:2181"}, content)
	if err != nil {
		return err
	}

	// BoltDB
	err = putData(store.BOLTDB, []string{"/tmp/test-traefik-certs-dumper.db"}, content)
	if err != nil {
		return err
	}

	return nil
}

func putData(backend store.Backend, addrs []string, content []byte) error {
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

	kvStore, err := valkeyrie.NewStore(backend, addrs, storeConfig)
	if err != nil {
		return err
	}

	if err := kvStore.Put(storeKey, content, nil); err != nil {
		return err
	}

	log.Printf("Successfully updated %s.\n", backend)
	return nil
}

func readFile(source string) ([]byte, error) {
	content, err := ioutil.ReadFile(filepath.Clean(source))
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
