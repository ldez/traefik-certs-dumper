package v2

import (
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/traefik/traefik/v2/pkg/provider/acme"
)

const (
	certsSubDir = "certs"
	keysSubDir  = "private"
)

// Dump Dumps data to certificates.
func Dump(data map[string]*acme.StoredData, baseConfig *dumper.BaseConfig) error {
	if baseConfig.Clean {
		err := cleanDir(baseConfig.DumpPath)
		if err != nil {
			return fmt.Errorf("folder cleaning failed: %w", err)
		}
	}

	if !baseConfig.DomainSubDir {
		if err := os.MkdirAll(filepath.Join(baseConfig.DumpPath, certsSubDir), 0755); err != nil {
			return fmt.Errorf("certs folder creation failure: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Join(baseConfig.DumpPath, keysSubDir), 0755); err != nil {
		return fmt.Errorf("keys folder creation failure: %w", err)
	}

	for _, store := range data {
		for _, cert := range store.Certificates {
			err := writeCert(baseConfig.DumpPath, cert.Certificate, baseConfig.CrtInfo, baseConfig.DomainSubDir)
			if err != nil {
				return fmt.Errorf("failed to write certificates: %w", err)
			}

			err = writeKey(baseConfig.DumpPath, cert.Certificate, baseConfig.KeyInfo, baseConfig.DomainSubDir)
			if err != nil {
				return fmt.Errorf("failed to write certificate keys: %w", err)
			}
		}

		if store.Account == nil {
			continue
		}

		privateKeyPem := extractPEMPrivateKey(store.Account)

		err := os.WriteFile(filepath.Join(baseConfig.DumpPath, keysSubDir, "letsencrypt"+baseConfig.KeyInfo.Ext), privateKeyPem, 0600)
		if err != nil {
			return fmt.Errorf("failed to write private key: %w", err)
		}
	}

	return nil
}

func writeCert(dumpPath string, cert acme.Certificate, info dumper.FileInfo, domainSubDir bool) error {
	certPath := filepath.Join(dumpPath, certsSubDir, safeName(cert.Domain.Main+info.Ext))
	if domainSubDir {
		certPath = filepath.Join(dumpPath, safeName(cert.Domain.Main), info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, safeName(cert.Domain.Main)), 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(certPath, cert.Certificate, 0666)
}

func writeKey(dumpPath string, cert acme.Certificate, info dumper.FileInfo, domainSubDir bool) error {
	keyPath := filepath.Join(dumpPath, keysSubDir, safeName(cert.Domain.Main+info.Ext))
	if domainSubDir {
		keyPath = filepath.Join(dumpPath, safeName(cert.Domain.Main), info.Name+info.Ext)
		if err := os.MkdirAll(filepath.Join(dumpPath, safeName(cert.Domain.Main)), 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(keyPath, cert.Key, 0600)
}

func extractPEMPrivateKey(account *acme.Account) []byte {
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

	dir, err := os.ReadDir(dumpPath)
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
