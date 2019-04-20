package dumper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"
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
	Domains     Domain
	Certificate *CertificateV1
}

// ConvertAccountV1ToV2 converts account information from version 1 to 2
func ConvertAccountV1ToV2(account *AccountV1) *StoredData {
	storedData := &StoredData{}
	storedData.Account = &Account{
		PrivateKey:   account.PrivateKey,
		Registration: account.Registration,
		Email:        account.Email,
		KeyType:      account.KeyType,
	}
	var certs []*Certificate
	for _, oldCert := range account.DomainsCertificate.Certs {
		certs = append(certs, &Certificate{
			Certificate: oldCert.Certificate.Certificate,
			Domain:      oldCert.Domains,
			Key:         oldCert.Certificate.PrivateKey,
		})
	}
	storedData.Certificates = certs
	return storedData
}

// Tree FIXME move
func Tree(root, indent string) error {
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}

	fmt.Println(fi.Name())
	if !fi.IsDir() {
		return nil
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string
	for _, fi := range fis {
		if fi.Name()[0] != '.' {
			names = append(names, fi.Name())
		}
	}

	for i, name := range names {
		add := "│  "
		if i == len(names)-1 {
			fmt.Printf(indent + "└──")
			add = "   "
		} else {
			fmt.Printf(indent + "├──")
		}

		if err := Tree(filepath.Join(root, name), indent+add); err != nil {
			return err
		}
	}

	return nil
}
