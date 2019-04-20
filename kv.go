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

const storeKey = "traefik/acme/account/object"

func getStoredDataFromGzip(value []byte) (*StoredData, error) {
	data := &StoredData{}

	r, err := gzip.NewReader(bytes.NewBuffer(value))
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
	case Consul:
		consul.Register()
		return store.CONSUL, nil
	case Etcd:
		etcdv3.Register()
		return store.ETCDV3, nil
	case Zookeeper:
		zookeeper.Register()
		return store.ZK, nil
	case BoldDB:
		boltdb.Register()
		return store.BOLTDB, nil
	default:
		return "", fmt.Errorf("no backend found for %v", backend)
	}
}

func loopKV(watch bool, kvStore store.Store, dataCh chan *StoredData, errCh chan error) {
	stopCh := make(<-chan struct{})
	events, err := kvStore.Watch(storeKey, stopCh, nil)
	if err != nil {
		errCh <- err
	}

	for {
		kvPair := <-events
		if kvPair == nil {
			errCh <- fmt.Errorf("could not fetch Key/Value pair for key %v", storeKey)
			return
		}
		dataCh <- extractStoredData(kvPair, errCh)
		if !watch {
			close(dataCh)
			close(errCh)
		}
	}
}

func extractStoredData(kvPair *store.KVPair, errCh chan error) *StoredData {
	storedData, err := getStoredDataFromGzip(kvPair.Value)
	if err != nil {
		errCh <- err
	}
	return storedData
}

func getSingleData(kvStore store.Store, dataCh chan *StoredData, errCh chan error) {
	kvPair, err := kvStore.Get(storeKey, nil)
	if err != nil {
		errCh <- err
		return
	}
	if kvPair == nil {
		errCh <- fmt.Errorf("could not fetch Key/Value pair for key %v", storeKey)
		return
	}

	dataCh <- extractStoredData(kvPair, errCh)
	close(dataCh)
	close(errCh)
}

func (b KVBackend) getStoredData(watch bool) (<-chan *StoredData, <-chan error) {
	dataCh := make(chan *StoredData)
	errCh := make(chan error)

	backend, err := register(b.Name)
	if err != nil {
		go func() {
			errCh <- err
		}()
		return dataCh, errCh
	}
	kvStore, err := valkeyrie.NewStore(
		backend,
		b.Client,
		b.Config,
	)

	if err != nil {
		go func() {
			errCh <- err
		}()
		return dataCh, errCh
	}

	if !watch {
		go getSingleData(kvStore, dataCh, errCh)
		return dataCh, errCh
	}

	go loopKV(watch, kvStore, dataCh, errCh)

	return dataCh, errCh

}
