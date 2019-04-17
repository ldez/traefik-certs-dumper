package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/zookeeper"
	"github.com/xenolf/lego/certcrypto"
	"github.com/xenolf/lego/registration"
)

const (
	certsSubDir = "certs"
	keysSubDir  = "private"
	storeKey    = "traefik/acme/account/object"
)

// StoredData represents the data managed by the Store
type StoredData struct {
	Account        *Account
	Certificates   []*Certificate
	HTTPChallenges map[string]map[string][]byte
	TLSChallenges  map[string]*Certificate
}

// Certificate is a struct which contains all data needed from an ACME certificate
type Certificate struct {
	Domain      Domain
	Certificate []byte
	Key         []byte
}

// Domain holds a domain name with SANs
type Domain struct {
	Main string
	SANs []string
}

// Account is used to store lets encrypt registration info
type Account struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   []byte
	KeyType      certcrypto.KeyType
}

type fileInfo struct {
	Name string
	Ext  string
}

// Backend represents a data source for ACME data
type Backend string

const (
	// FILE backend
	FILE Backend = "file"
	// CONSUL backend
	CONSUL Backend = "consul"
	// ETCD backend
	ETCD Backend = "etcd"
	// ZK backend
	ZK Backend = "zk"
	// BOLTDB backend
	BOLTDB Backend = "boltdb"
)

type dumpConfig struct {
	Path         string
	CertInfo     fileInfo
	KeyInfo      fileInfo
	DomainSubDir bool
	Watch        bool
}

func getAcmeDataFromJSONFile(file string) (*StoredData, error) {
	data := &StoredData{}

	f, err := os.Open(file)
	if err != nil {
		return data, err
	}
	fmt.Println(data)

	if err = json.NewDecoder(f).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

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

func loop(config *dumpConfig, backend Backend) error {

	// TODO change to env parameter
	client := "localhost:8500"

	var storeBackend store.Backend
	switch backend {
	case CONSUL:
		storeBackend = store.CONSUL
		consul.Register()
	case ETCD:
		storeBackend = store.ETCDV3
		etcdv3.Register()
	case ZK:
		storeBackend = store.ZK
		zookeeper.Register()
	case BOLTDB:
		storeBackend = store.BOLTDB
		boltdb.Register()
	}
	kvstore, err := valkeyrie.NewStore(
		storeBackend,
		[]string{client},
		&store.Config{},
	)
	if err != nil {
		return err
	}

	stopCh := make(<-chan struct{})
	events, _ := kvstore.Watch(storeKey, stopCh, nil)
	for {
		select {
		case kvpair := <-events:
			storedData, err := getStoredDataFromGzip(kvpair.Value)
			if err != nil {
				return err
			}
			if err := dump(config, storedData); err != nil {
				return err
			}
			if !config.Watch {
				return nil
			}
		}
	}
}

func dump(config *dumpConfig, data *StoredData) error {

	if err := os.RemoveAll(config.Path); err != nil {
		return err
	}

	if !config.DomainSubDir {
		if err := os.MkdirAll(filepath.Join(config.Path, certsSubDir), 0755); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(config.Path, keysSubDir), 0755); err != nil {
		return err
	}

	privateKeyPem := extractPEMPrivateKey(data.Account)
	err := ioutil.WriteFile(filepath.Join(config.Path, keysSubDir, "letsencrypt"+config.KeyInfo.Ext), privateKeyPem, 0666)
	if err != nil {
		return err
	}

	for _, cert := range data.Certificates {
		err := writeCert(config.Path, cert, config.CertInfo, config.DomainSubDir)
		if err != nil {
			return err
		}

		err = writeKey(config.Path, cert, config.KeyInfo, config.DomainSubDir)
		if err != nil {
			return err
		}
	}

	if err := tree(config.Path, ""); err != nil {
		return err
	}

	return nil
}

func writeCert(dumpPath string, cert *Certificate, info fileInfo, domainSubDir bool) error {
	certPath := filepath.Join(dumpPath, certsSubDir, cert.Domain.Main+info.Ext)
	if domainSubDir {
		certPath = filepath.Join(dumpPath, cert.Domain.Main, info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, cert.Domain.Main), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(certPath, cert.Certificate, 0666)
}

func writeKey(dumpPath string, cert *Certificate, info fileInfo, domainSubDir bool) error {
	keyPath := filepath.Join(dumpPath, keysSubDir, cert.Domain.Main+info.Ext)
	if domainSubDir {
		keyPath = filepath.Join(dumpPath, cert.Domain.Main, info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, cert.Domain.Main), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(keyPath, cert.Key, 0666)
}

func extractPEMPrivateKey(account *Account) []byte {
	var block *pem.Block
	switch account.KeyType {
	case certcrypto.RSA2048, certcrypto.RSA4096, certcrypto.RSA8192:
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: account.PrivateKey,
		}
	case certcrypto.EC256, certcrypto.EC384:
		block = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: account.PrivateKey,
		}
	default:
		panic("unsupported key type")
	}

	return pem.EncodeToMemory(block)
}

func tree(root, indent string) error {
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}

	fmt.Println(fi.Name())
	if !fi.IsDir() {
		return nil
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string
	for _, fi := range fis {
		if fi.Name()[0] != '.' {
			names = append(names, fi.Name())
		}
	}

	for i, name := range names {
		add := "│  "
		if i == len(names)-1 {
			fmt.Printf(indent + "└──")
			add = "   "
		} else {
			fmt.Printf(indent + "├──")
		}

		if err := tree(filepath.Join(root, name), indent+add); err != nil {
			return err
		}
	}

	return nil
}
