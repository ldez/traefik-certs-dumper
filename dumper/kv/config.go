package kv

import "github.com/abronan/valkeyrie/store"

// Config FIXME
type Config struct {
	Backend   store.Backend
	Prefix    string
	Endpoints []string
	Options   *store.Config
}
