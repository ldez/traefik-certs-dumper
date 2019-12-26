package kv

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	v1 "github.com/ldez/traefik-certs-dumper/v2/dumper/v1"
	"github.com/ldez/traefik-certs-dumper/v2/hook"
)

// DefaultStoreKeySuffix is the default suffix/storage.
const DefaultStoreKeySuffix = "/acme/account/object"

// Dump Dumps KV content to certificates.
func Dump(config *Config, baseConfig *dumper.BaseConfig) error {
	kvStore, err := valkeyrie.NewStore(config.Backend, config.Endpoints, config.Options)
	if err != nil {
		return fmt.Errorf("unable to create client of the store: %w", err)
	}

	storeKey := config.Prefix + config.Suffix

	if baseConfig.Watch {
		return watch(kvStore, storeKey, baseConfig)
	}

	pair, err := kvStore.Get(storeKey, nil)
	if err != nil {
		return fmt.Errorf("unable to retrieve %s value: %w", storeKey, err)
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

		if isDebug() {
			log.Println("Dumped new certificate data.")
		}

		hook.Exec(baseConfig.Hook)
	}
}

func dumpPair(pair *store.KVPair, baseConfig *dumper.BaseConfig) error {
	data, err := getStoredDataFromGzip(pair)
	if err != nil {
		return err
	}

	return v1.Dump(data, baseConfig)
}

func getStoredDataFromGzip(pair *store.KVPair) (*v1.StoredData, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(pair.Value))
	if err != nil {
		return nil, fmt.Errorf("fail to create GZip reader: %w", err)
	}

	acmeData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read the pair content: %w", err)
	}

	account := &AccountOld{}
	if err := json.Unmarshal(acmeData, &account); err != nil {
		return nil, fmt.Errorf("unable marshal AccountOld: %w", err)
	}

	return convertOldAccount(account), nil
}

func isDebug() bool {
	return strings.EqualFold(os.Getenv("TCD_DEBUG"), "true")
}
