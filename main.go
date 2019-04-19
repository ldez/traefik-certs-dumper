package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "traefik-certs-dumper",
		Short:   "Dump Let's Encrypt certificates from Traefik",
		Long:    `Dump ACME data from Traefik of different storage backends to certificates.`,
		Version: version,
	}

	config := &Config{}

	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump Let's Encrypt certificates from Traefik",
		Long:  `Dump ACME data from Traefik of different storage backends to certificates.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			source := cmd.Flag("source").Value.String()
			sourceFile := cmd.Flag("source.file").Value.String()
			watch, _ := strconv.ParseBool(cmd.Flag("watch").Value.String())

			switch source {
			case FILE:
				if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
					return fmt.Errorf("--source.file (%q) does not exist", sourceFile)
				}
			case BOLTDB:
				if watch {
					return fmt.Errorf("--watch=true is not supported for boltdb")
				}
			case CONSUL:
			case ETCD:
			case ZOOKEEPER:
			default:
				return fmt.Errorf("--source (%q) is not allowed, use one of 'file', 'consul', 'etcd', 'zookeeper', 'boltdb'", source)
			}

			crtExt := cmd.Flag("crt-ext").Value.String()
			keyExt := cmd.Flag("key-ext").Value.String()

			subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
			if !subDir {
				if crtExt == keyExt {
					return fmt.Errorf("--crt-ext (%q) and --key-ext (%q) are identical, in this case --domain-subdir is required", crtExt, keyExt)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {

			source := cmd.Flag("source").Value.String()
			acmeFile := cmd.Flag("source.file").Value.String()

			endpoints := strings.Split(cmd.Flag("source.kv.endpoints").Value.String(), ",")

			storeConfig := &store.Config{}

			timeout, _ := strconv.Atoi(cmd.Flag("source.kv.connection-timeout").Value.String())
			storeConfig.ConnectionTimeout = time.Second * time.Duration(timeout)
			storeConfig.Username = cmd.Flag("source.kv.username").Value.String()
			storeConfig.Password = cmd.Flag("source.kv.password").Value.String()

			enableTLS, _ := strconv.ParseBool(cmd.Flag("source.kv.tls.enable").Value.String())

			if enableTLS {
				tlsConfig := &tls.Config{}
				insecureSkipVerify, _ := strconv.ParseBool(cmd.Flag("source.kv.tls.insecureskipverify").Value.String())
				tlsConfig.InsecureSkipVerify = insecureSkipVerify
				if cmd.Flag("source.kv.tls.ca-cert-file").Value.String() != "" {
					caFile := cmd.Flag("source.kv.tls.ca-cert-file").Value.String()
					caCert, err := ioutil.ReadFile(caFile)
					if err != nil {
						log.Fatal(err)
					}
					roots := x509.NewCertPool()
					ok := roots.AppendCertsFromPEM(caCert)
					if !ok {
						log.Fatalf("failed to parse root certificate")
					}
					tlsConfig.RootCAs = roots
				}
				storeConfig.TLS = tlsConfig
			}

			// Special parameters for etcd
			timeout, _ = strconv.Atoi(cmd.Flag("source.kv.etcd.sync-period").Value.String())
			storeConfig.SyncPeriod = time.Second * time.Duration(timeout)
			// Special parameters for boltdb
			persistConnection, _ := strconv.ParseBool(cmd.Flag("source.kv.boltdb.persist-connection").Value.String())
			storeConfig.PersistConnection = persistConnection
			storeConfig.Bucket = cmd.Flag("source.kv.boltdb.bucket").Value.String()
			// Special parameters for consul
			storeConfig.Token = cmd.Flag("source.kv.consul.token").Value.String()

			switch source {
			case "file":
				config.BackendConfig = FileBackend{
					Name: FILE,
					Path: acmeFile,
				}
			case "consul":
				fallthrough
			case "etcd":
				fallthrough
			case "zookeeper":
				fallthrough
			case "boltdb":
				fallthrough
			default:
				config.BackendConfig = KVBackend{
					Name:   source,
					Client: endpoints,
					Config: storeConfig,
				}
			}

			config.Path = cmd.Flag("dest").Value.String()

			config.CertInfo = fileInfo{
				Name: cmd.Flag("crt-name").Value.String(),
				Ext:  cmd.Flag("crt-ext").Value.String(),
			}

			config.KeyInfo = fileInfo{
				Name: cmd.Flag("key-name").Value.String(),
				Ext:  cmd.Flag("key-ext").Value.String(),
			}

			config.DomainSubDir, _ = strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
			config.Watch, _ = strconv.ParseBool(cmd.Flag("watch").Value.String())

			if err := run(config); err != nil {
				fmt.Println(err)
			}

			return nil
		},
	}

	dumpCmd.Flags().String("source", "file", "Source type, one of 'file', 'consul', 'etcd', 'zookeeper', 'boltdb'. Options for each source type are prefixed with `source.<type>.`")
	dumpCmd.Flags().String("source.file", "./acme.json", "Path to 'acme.json' for file source.")

	// Generic parameters for Key/Value backends
	dumpCmd.Flags().String("source.kv.endpoints", "localhost:8500", "Comma seperated list of endpoints.")
	dumpCmd.Flags().Int("source.kv.connection-timeout", 0, "Connection timeout in seconds.")
	dumpCmd.Flags().String("source.kv.password", "", "Password for connection.")
	dumpCmd.Flags().String("source.kv.username", "", "Username for connection.")
	dumpCmd.Flags().Bool("source.kv.tls.enable", false, "Enable TLS encryption.")
	dumpCmd.Flags().Bool("source.kv.tls.insecureskipverify", false, "Trust unverified certificates if TLS is enabled.")
	dumpCmd.Flags().String("source.kv.tls.ca-cert-file", "", "Root CA file for certificate verification if TLS is enabled.")
	// Special parameters for etcd
	dumpCmd.Flags().Int("source.kv.etcd.sync-period", 0, "Sync period for etcd in seconds.")
	// Special parameters for boltdb
	dumpCmd.Flags().Bool("source.kv.boltdb.persist-connection", false, "Persist connection for boltdb.")
	dumpCmd.Flags().String("source.kv.boltdb.bucket", "traefik", "Bucket for boltdb.")
	// Special parameters for consul
	dumpCmd.Flags().String("source.kv.consul.token", "", "Token for consul.")

	dumpCmd.Flags().Bool("watch", false, "Enable watching changes.")
	dumpCmd.Flags().String("dest", "./dump", "Path to store the dump content.")
	dumpCmd.Flags().String("crt-ext", ".crt", "The file extension of the generated certificates.")
	dumpCmd.Flags().String("crt-name", "certificate", "The file name (without extension) of the generated certificates.")
	dumpCmd.Flags().String("key-ext", ".key", "The file extension of the generated private keys.")
	dumpCmd.Flags().String("key-name", "privatekey", "The file name (without extension) of the generated private keys.")
	dumpCmd.Flags().Bool("domain-subdir", false, "Use domain as sub-directory.")
	rootCmd.AddCommand(dumpCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Run: func(_ *cobra.Command, _ []string) {
			displayVersion(rootCmd.Name())
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
