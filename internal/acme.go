package internal

import "github.com/go-acme/lego/v5/acme"

type Resource struct {
	Body acme.Account `json:"body"`
	URI  string       `json:"uri,omitempty"`
}
