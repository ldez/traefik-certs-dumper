package main

// Certificates Data Sources.
const (
	File      = "file"
	Consul    = "consul"
	Etcd      = "etcd"
	Zookeeper = "zookeeper"
	BoldDB    = "boltdb"
)

// Config represents a configuration for dumping certificates
type Config struct {
	Path          string
	CertInfo      fileInfo
	KeyInfo       fileInfo
	DomainSubDir  bool
	Watch         bool
	BackendConfig interface{}
}

// Backend represents an object storage of ACME data
type Backend interface {
	getStoredData(watch bool) (<-chan *StoredData, <-chan error)
}

func run(config *Config) error {
	data, errors := config.BackendConfig.(Backend).getStoredData(config.Watch)
	for {
		select {
		case err := <-errors:
			return err
		case acmeData, ok := <-data:
			if !ok {
				return nil
			}
			if err := dump(config, acmeData); err != nil {
				return err
			}
		}
	}
}
