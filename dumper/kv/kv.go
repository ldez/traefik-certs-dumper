package kv

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/ldez/traefik-certs-dumper/dumper"
)

const storeKey = "/acme/account/object"

// Dump FIXME
func Dump(config *Config, baseConfig *dumper.BaseConfig) error {
	kvStore, err := valkeyrie.NewStore(config.Backend, config.Endpoints, config.Options)
	if err != nil {
		return err
	}

	pair, err := kvStore.Get(config.Prefix+storeKey, nil)
	if err != nil {
		return err
	}

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

	storedData := &dumper.StoredData{}
	if err := json.Unmarshal(acmeData, &storedData); err != nil {
		return data, err
	}

	return storedData, nil
}
