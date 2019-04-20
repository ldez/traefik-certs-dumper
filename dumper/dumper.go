package dumper

import (
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/certcrypto"
)

const (
	certsSubDir = "certs"
	keysSubDir  = "private"
)

// FileInfo File information.
type FileInfo struct {
	Name string
	Ext  string
}

// Dump Dumps data to certificates.
func Dump(data *StoredData, baseConfig *BaseConfig) error {

	if baseConfig.Clean {
		if err := os.RemoveAll(baseConfig.DumpPath); err != nil {
			return err
		}
	}

	if !baseConfig.DomainSubDir {
		if err := os.MkdirAll(filepath.Join(baseConfig.DumpPath, certsSubDir), 0755); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(baseConfig.DumpPath, keysSubDir), 0755); err != nil {
		return err
	}

	privateKeyPem := extractPEMPrivateKey(data.Account)
	err := ioutil.WriteFile(filepath.Join(baseConfig.DumpPath, keysSubDir, "letsencrypt"+baseConfig.KeyInfo.Ext), privateKeyPem, 0666)
	if err != nil {
		return err
	}

	for _, cert := range data.Certificates {
		err := writeCert(baseConfig.DumpPath, cert, baseConfig.CrtInfo, baseConfig.DomainSubDir)
		if err != nil {
			return err
		}

		err = writeKey(baseConfig.DumpPath, cert, baseConfig.KeyInfo, baseConfig.DomainSubDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeCert(dumpPath string, cert *Certificate, info FileInfo, domainSubDir bool) error {
	certPath := filepath.Join(dumpPath, certsSubDir, cert.Domain.Main+info.Ext)
	if domainSubDir {
		certPath = filepath.Join(dumpPath, cert.Domain.Main, info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, cert.Domain.Main), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(certPath, cert.Certificate, 0666)
}

func writeKey(dumpPath string, cert *Certificate, info FileInfo, domainSubDir bool) error {
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
