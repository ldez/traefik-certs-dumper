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

	"github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/watch"
	"github.com/xenolf/lego/certcrypto"
	"github.com/xenolf/lego/registration"
)

const (
	certsSubDir = "certs"
	keysSubDir  = "private"
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

func dump(data StoredData, dumpPath string, crtInfo, keyInfo fileInfo, domainSubDir bool) error {

	if err := os.RemoveAll(dumpPath); err != nil {
		return err
	}

	if !domainSubDir {
		if err := os.MkdirAll(filepath.Join(dumpPath, certsSubDir), 0755); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(dumpPath, keysSubDir), 0755); err != nil {
		return err
	}

	privateKeyPem := extractPEMPrivateKey(data.Account)
	err := ioutil.WriteFile(filepath.Join(dumpPath, keysSubDir, "letsencrypt"+keyInfo.Ext), privateKeyPem, 0666)
	if err != nil {
		return err
	}

	for _, cert := range data.Certificates {
		err := writeCert(dumpPath, cert, crtInfo, domainSubDir)
		if err != nil {
			return err
		}

		err = writeKey(dumpPath, cert, keyInfo, domainSubDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpFile(acmeFile string, dumpPath string, crtInfo, keyInfo fileInfo, domainSubDir bool) error {
	f, err := os.Open(acmeFile)
	if err != nil {
		return err
	}

	data := StoredData{}
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	return dump(data, dumpPath, crtInfo, keyInfo, domainSubDir)
}

func dumpConsul(watch bool, dumpPath string, crtInfo, keyInfo fileInfo, domainSubDir bool) {

	params := map[string]interface{}{
		"type": "key",
		"key":  "traefik/acme/account/object",
	}
	plan, _ := consulwatch.Parse(params)
	plan.Handler = func(idx uint64, data interface{}) {

		// TODO is here a better way?
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(data)
		kvpair := api.KVPair{}
		json.Unmarshal(buf.Bytes(), &kvpair)

		r, err := gzip.NewReader(bytes.NewBuffer(kvpair.Value))
		defer r.Close()
		if err != nil {
			fmt.Printf("[ERR] %s", err)
		}

		acmeData, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Printf("[ERR] %s", err)
		}

		storedData := StoredData{}
		json.Unmarshal(acmeData, &storedData)

		dump(storedData, dumpPath, crtInfo, keyInfo, domainSubDir)
		if err := tree(dumpPath, ""); err != nil {
			fmt.Printf("[ERR] %s", err)
		}
		if !watch {
			plan.Stop()
		}
	}

	config := api.DefaultConfig()
	fmt.Println("Start watching consul...")
	plan.Run(config.Address)
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
