package main

import (
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xenolf/lego/certcrypto"
	"github.com/xenolf/lego/registration"
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

func dump(acmeFile, dumpPath string, crtExt, keyExt string) error {
	f, err := os.Open(acmeFile)
	if err != nil {
		return err
	}

	data := StoredData{}
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	if err = os.RemoveAll(dumpPath); err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(dumpPath, "certs"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(dumpPath, "private"), 0755)
	if err != nil {
		return err
	}

	privateKeyPem := extractPEMPrivateKey(data.Account)
	err = ioutil.WriteFile(filepath.Join(dumpPath, "private", "letsencrypt"+keyExt), privateKeyPem, 0666)
	if err != nil {
		return err
	}

	for _, cert := range data.Certificates {
		err = ioutil.WriteFile(filepath.Join(dumpPath, "private", cert.Domain.Main+keyExt), cert.Key, 0666)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(dumpPath, "certs", cert.Domain.Main+crtExt), cert.Certificate, 0666)
		if err != nil {
			return err
		}
	}

	return nil
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
