package main

import (
	"log"

	"github.com/ldez/traefik-certs-dumper/v2/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
