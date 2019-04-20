package main

import (
	"log"

	"github.com/ldez/traefik-certs-dumper/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
