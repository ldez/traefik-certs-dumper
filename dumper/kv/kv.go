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

// FIXME prefix
const storeKey = "traefik/acme/account/object"

// BaseConfig FIXME
type BaseConfig struct {
	Backend   store.Backend
	Endpoints []string
	Options   *store.Config
}

// Dump FIXME
func Dump(config *BaseConfig, dumpPath string, crtInfo, keyInfo dumper.FileInfo, domainSubDir bool) error {
	kvStore, err := valkeyrie.NewStore(config.Backend, config.Endpoints, config.Options)
	if err != nil {
		return err
	}

	pair, err := kvStore.Get(storeKey, nil)
	if err != nil {
		return err
	}

	data, err := getStoredDataFromGzip(pair)
	if err != nil {
		return err
	}

	return dumper.Dump(data, dumpPath, crtInfo, keyInfo, domainSubDir)
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
