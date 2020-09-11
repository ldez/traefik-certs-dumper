package v1

import (
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/registration"
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
	Registration *registration.Resource
	PrivateKey   []byte
	KeyType      certcrypto.KeyType
}
