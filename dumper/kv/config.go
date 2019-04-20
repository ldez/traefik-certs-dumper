package kv

import "github.com/abronan/valkeyrie/store"

// Config FIXME
type Config struct {
	Backend   store.Backend
	Endpoints []string
	Options   *store.Config
}
