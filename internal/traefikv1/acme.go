package traefikv1

import (
	"github.com/ldez/traefik-certs-dumper/v2/internal"
)

// StoredData represents the data managed by the Store.
type StoredData struct {
	Account        *Account
	Certificates   []*Certificate
	HTTPChallenges map[string]map[string][]byte
	TLSChallenges  map[string]*Certificate
}

// Certificate is a struct which contains all data needed from an ACME certificate.
type Certificate struct {
	Domain      Domain
	Certificate []byte
	Key         []byte
}

// Domain holds a domain name with SANs.
type Domain struct {
	Main string
	SANs []string
}

// Account is used to store lets encrypt registration info.
type Account struct {
	Email        string
	Registration *internal.Resource
	PrivateKey   []byte
}
