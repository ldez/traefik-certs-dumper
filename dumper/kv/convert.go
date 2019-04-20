package kv

import (
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
)

// CertificateV1 is used to store certificate info
type CertificateV1 struct {
	Domain        string
	CertURL       string
	CertStableURL string
	PrivateKey    []byte
	Certificate   []byte
}

// AccountV1 is used to store lets encrypt registration info
type AccountV1 struct {
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
	Domains     dumper.Domain
	Certificate *CertificateV1
}

// convertAccountV1ToV2 converts account information from version 1 to 2
func convertAccountV1ToV2(account *AccountV1) *dumper.StoredData {
	storedData := &dumper.StoredData{}
	storedData.Account = &dumper.Account{
		PrivateKey:   account.PrivateKey,
		Registration: account.Registration,
		Email:        account.Email,
		KeyType:      account.KeyType,
	}
	var certs []*dumper.Certificate
	for _, oldCert := range account.DomainsCertificate.Certs {
		certs = append(certs, &dumper.Certificate{
			Certificate: oldCert.Certificate.Certificate,
			Domain:      oldCert.Domains,
			Key:         oldCert.Certificate.PrivateKey,
		})
	}
	storedData.Certificates = certs
	return storedData
}
