package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/zookeeper"
)

const (
	storeKey = "traefik/acme/account/object"
)

func getStoredDataFromGzip(value []byte) (*StoredData, error) {
	data := &StoredData{}

	r, err := gzip.NewReader(bytes.NewBuffer(value))
	defer r.Close()
	if err != nil {
		return data, err
	}

	acmeData, err := ioutil.ReadAll(r)
	if err != nil {
		return data, err
	}

	storedData := &StoredData{}
	if err := json.Unmarshal(acmeData, &storedData); err != nil {
		return data, err
	}

	return storedData, nil
}

// KVBackend represents a Key/Value pair backend
type KVBackend struct {
	Name   string
	Client []string
	Config *store.Config
}

func register(backend string) (store.Backend, error) {
	switch backend {
	case CONSUL:
		consul.Register()
		return store.CONSUL, nil
	case ETCD:
		etcdv3.Register()
		return store.ETCDV3, nil
	case ZK:
		zookeeper.Register()
		return store.ZK, nil
	case BOLTDB:
		boltdb.Register()
		return store.BOLTDB, nil
	default:
		return "", fmt.Errorf("No backend found for %v", backend)
	}
}

func (b KVBackend) loop(watch bool) (<-chan *StoredData, <-chan error) {

	dataCh := make(chan *StoredData)
	errors := make(chan error)

	backend, err := register(b.Name)
	if err != nil {
		errors <- err
	}

	kvstore, err := valkeyrie.NewStore(
		backend,
		b.Client,
		b.Config,
	)
	if err != nil {
		errors <- err
	}

	go func() {
		stopCh := make(<-chan struct{})
		events, _ := kvstore.Watch(storeKey, stopCh, nil)
		for {
			select {
			case kvpair := <-events:
				storedData, err := getStoredDataFromGzip(kvpair.Value)
				if err != nil {
					errors <- err
				}
				dataCh <- storedData
			}
			if !watch {
				close(dataCh)
				close(errors)
			}
		}
	}()

	return dataCh, errors

}
