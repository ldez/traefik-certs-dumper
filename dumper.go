package main

const (
	// FILE backend
	FILE string = "file"
	// CONSUL backend
	CONSUL string = "consul"
	// ETCD backend
	ETCD string = "etcd"
	// ZOOKEEPER backend
	ZOOKEEPER string = "zookeeper"
	// BOLTDB backend
	BOLTDB string = "boltdb"
)

// Config represents a configuration for dumping cerificates
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
