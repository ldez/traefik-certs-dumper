package dumper

import (
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v3/certcrypto"
)

const (
	certsSubDir = "certs"
	keysSubDir  = "private"
)

// Dump Dumps data to certificates.
func Dump(data *StoredData, baseConfig *BaseConfig) error {
	if baseConfig.Clean {
		err := cleanDir(baseConfig.DumpPath)
		if err != nil {
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
	err := ioutil.WriteFile(filepath.Join(baseConfig.DumpPath, keysSubDir, "letsencrypt"+baseConfig.KeyInfo.Ext), privateKeyPem, 0600)
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
	certPath := filepath.Join(dumpPath, certsSubDir, safeName(cert.Domain.Main+info.Ext))
	if domainSubDir {
		certPath = filepath.Join(dumpPath, safeName(cert.Domain.Main), info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, safeName(cert.Domain.Main)), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(certPath, cert.Certificate, 0666)
}

func writeKey(dumpPath string, cert *Certificate, info FileInfo, domainSubDir bool) error {
	keyPath := filepath.Join(dumpPath, keysSubDir, safeName(cert.Domain.Main+info.Ext))
	if domainSubDir {
		keyPath = filepath.Join(dumpPath, safeName(cert.Domain.Main), info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, safeName(cert.Domain.Main)), 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(keyPath, cert.Key, 0600)
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

func cleanDir(dumpPath string) error {
	_, errExists := os.Stat(dumpPath)
	if os.IsNotExist(errExists) {
		return nil
	}

	if errExists != nil {
		return errExists
	}

	dir, err := ioutil.ReadDir(dumpPath)
	if err != nil {
		return err
	}

	for _, f := range dir {
		if err := os.RemoveAll(filepath.Join(dumpPath, f.Name())); err != nil {
			return err
		}
	}

	return nil
}
