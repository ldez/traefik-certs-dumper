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

func dump(acmeFile, dumpPath string, crtInfo, keyInfo fileInfo, domainSubDir bool) error {
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

	if !domainSubDir {
		if err = os.MkdirAll(filepath.Join(dumpPath, certsSubDir), 0755); err != nil {
			return err
		}
	}

	if err = os.MkdirAll(filepath.Join(dumpPath, keysSubDir), 0755); err != nil {
		return err
	}

	privateKeyPem := extractPEMPrivateKey(data.Account)
	err = ioutil.WriteFile(filepath.Join(dumpPath, keysSubDir, "letsencrypt"+keyInfo.Ext), privateKeyPem, 0666)
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

func writeCert(dumpPath string, cert *Certificate, info fileInfo, domainSubDir bool) error {
	certPath := filepath.Join(dumpPath, keysSubDir, cert.Domain.Main+info.Ext)
	if domainSubDir {
		certPath = filepath.Join(dumpPath, cert.Domain.Main, info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, cert.Domain.Main), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(certPath, cert.Certificate, 0666)
}

func writeKey(dumpPath string, cert *Certificate, info fileInfo, domainSubDir bool) error {
	keyPath := filepath.Join(dumpPath, certsSubDir, cert.Domain.Main+info.Ext)
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
