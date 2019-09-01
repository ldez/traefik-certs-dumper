package kv

import (
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/registration"
	v1 "github.com/ldez/traefik-certs-dumper/v2/dumper/v1"
)

// CertificateOld is used to store certificate info
type CertificateOld struct {
	Domain        string
	CertURL       string
	CertStableURL string
	PrivateKey    []byte
	Certificate   []byte
}

// AccountOld is used to store lets encrypt registration info
type AccountOld struct {
	Email              string
	Registration       *registration.Resource
	PrivateKey         []byte
	KeyType            certcrypto.KeyType
	DomainsCertificate DomainsCertificates
	ChallengeCerts     map[string]*ChallengeCert
	HTTPChallenge      map[string]map[string][]byte
}

// DomainsCertificates stores a certificate for multiple domains
type DomainsCertificates struct {
	Certs []*DomainsCertificate
}

// ChallengeCert stores a challenge certificate
type ChallengeCert struct {
	Certificate []byte
	PrivateKey  []byte
}

// DomainsCertificate contains a certificate for multiple domains
type DomainsCertificate struct {
	Domains     v1.Domain
	Certificate *CertificateOld
}

// convertOldAccount converts account information from old account format.
func convertOldAccount(account *AccountOld) *v1.StoredData {
	storedData := &v1.StoredData{}
	storedData.Account = &v1.Account{
		PrivateKey:   account.PrivateKey,
		Registration: account.Registration,
		Email:        account.Email,
		KeyType:      account.KeyType,
	}
	var certs []*v1.Certificate
	for _, oldCert := range account.DomainsCertificate.Certs {
		certs = append(certs, &v1.Certificate{
			Certificate: oldCert.Certificate.Certificate,
			Domain:      oldCert.Domains,
			Key:         oldCert.Certificate.PrivateKey,
		})
	}
	storedData.Certificates = certs
	return storedData
}
