package kv

import (
	"github.com/kvtools/valkeyrie"
)

// Config KV configuration.
type Config struct {
	StoreName string
	Prefix    string
	Suffix    string
	Endpoints []string
	Options   valkeyrie.Config
}
