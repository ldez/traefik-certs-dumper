package kv

import "github.com/abronan/valkeyrie/store"

// Config KV configuration.
type Config struct {
	Backend   store.Backend
	Prefix    string
	Suffix    string
	Endpoints []string
	Options   *store.Config
}
