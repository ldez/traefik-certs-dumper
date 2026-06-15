package kv

import (
	"github.com/ldez/traefik-certs-dumper/v2/internal"
	"github.com/ldez/traefik-certs-dumper/v2/internal/traefikv1"
)

// CertificateOld is used to store certificate info.
type CertificateOld struct {
	Domain        string
	CertURL       string
	CertStableURL string
	PrivateKey    []byte
	Certificate   []byte
}

// AccountOld is used to store lets encrypt registration info.
type AccountOld struct {
	Email              string
	Registration       *internal.Resource
	PrivateKey         []byte
	DomainsCertificate DomainsCertificates
	ChallengeCerts     map[string]*ChallengeCert
	HTTPChallenge      map[string]map[string][]byte
}

// DomainsCertificates stores a certificate for multiple domains.
type DomainsCertificates struct {
	Certs []*DomainsCertificate
}

// ChallengeCert stores a challenge certificate.
type ChallengeCert struct {
	Certificate []byte
	PrivateKey  []byte
}

// DomainsCertificate contains a certificate for multiple domains.
type DomainsCertificate struct {
	Domains     traefikv1.Domain
	Certificate *CertificateOld
}

// convertOldAccount converts account information from old account format.
func convertOldAccount(account *AccountOld) *traefikv1.StoredData {
	storedData := &traefikv1.StoredData{
		Account: &traefikv1.Account{
			PrivateKey:   account.PrivateKey,
			Registration: account.Registration,
			Email:        account.Email,
		},
	}

	var certs []*traefikv1.Certificate
	for _, oldCert := range account.DomainsCertificate.Certs {
		certs = append(certs, &traefikv1.Certificate{
			Certificate: oldCert.Certificate.Certificate,
			Domain:      oldCert.Domains,
			Key:         oldCert.Certificate.PrivateKey,
		})
	}

	storedData.Certificates = certs

	return storedData
}
