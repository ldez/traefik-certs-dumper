package kv

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/ldez/traefik-certs-dumper/dumper"
)

const storeKeySuffix = "/acme/account/object"

// Dump FIXME
func Dump(config *Config, baseConfig *dumper.BaseConfig) error {
	kvStore, err := valkeyrie.NewStore(config.Backend, config.Endpoints, config.Options)
	if err != nil {
		return err
	}

	storeKey := config.Prefix + storeKeySuffix

	if config.Watch {
		return watch(kvStore, storeKey, baseConfig)
	}

	pair, err := kvStore.Get(storeKey, nil)
	if err != nil {
		return err
	}

	return dumpPair(pair, baseConfig)
}

func watch(kvStore store.Store, storeKey string, baseConfig *dumper.BaseConfig) error {
	stopCh := make(<-chan struct{})

	pairs, err := kvStore.Watch(storeKey, stopCh, nil)
	if err != nil {
		return err
	}

	for {
		pair := <-pairs
		if pair == nil {
			return fmt.Errorf("could not fetch Key/Value pair for key %v", storeKey)
		}

		err = dumpPair(pair, baseConfig)
		if err != nil {
			return err
		}
	}
}

func dumpPair(pair *store.KVPair, baseConfig *dumper.BaseConfig) error {
	data, err := getStoredDataFromGzip(pair)
	if err != nil {
		return err
	}

	return dumper.Dump(data, baseConfig)
}

func getStoredDataFromGzip(pair *store.KVPair) (*dumper.StoredData, error) {
	data := &dumper.StoredData{}

	reader, err := gzip.NewReader(bytes.NewBuffer(pair.Value))
	if err != nil {
		return data, err
	}

	acmeData, err := ioutil.ReadAll(reader)
	if err != nil {
		return data, err
	}

	account := &dumper.AccountV1{}
	if err := json.Unmarshal(acmeData, &account); err != nil {
		return data, err
	}

	return dumper.ConvertAccountV1ToV2(account), nil
}
