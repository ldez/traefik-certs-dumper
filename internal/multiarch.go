package main

import (
	"flag"
	"log"
	"os"
)

type buildOption struct {
	OS      string `json:"os"`
	GoARCH  string `json:"go_arch"`
	GoARM   string `json:"go_arm,omitempty"`
	Variant string `json:"variant,omitempty"`
}

func main() {
	log.SetFlags(log.Lshortfile)

	imageName := flag.String("image-name", "ldez/traefik-certs-dumper", "")
	version := flag.String("version", "", "")
	baseImageName := flag.String("base-image-name", "alpine:3.9", "")
	dryRun := flag.Bool("dry-run", true, "")

	flag.Parse()

	require("image-name", imageName)
	require("version", version)
	require("base-image-name", baseImageName)

	_, travisTag := os.LookupEnv("TRAVIS_TAG")
	if !travisTag {
		log.Println("Skipping deploy")
		os.Exit(0)
	}

	targets := []string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"}

	publisher, err := newPublisher(*imageName, *version, *baseImageName, targets)
	if err != nil {
		log.Fatal(err)
	}

	err = publisher.execute(*dryRun)
	if err != nil {
		log.Fatal(err)
	}
}

func require(fieldName string, field *string) {
	if field == nil || *field == "" {
		log.Fatalf("%s is required", fieldName)
	}
}
