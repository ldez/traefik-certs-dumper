package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abronan/valkeyrie/store"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "traefik-certs-dumper",
		Short:   "Dump Let's Encrypt certificates from Traefik",
		Long:    `Dump the content of the "acme.json" file from Traefik to certificates.`,
		Version: version,
	}

	config := &Config{}

	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump Let's Encrypt certificates from Traefik",
		Long:  `Dump the content of the "acme.json" file from Traefik to certificates.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			source := cmd.Flag("source").Value.String()
			sourceFile := cmd.Flag("source-file").Value.String()
			if source == "file" {
				if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
					return fmt.Errorf("--source-file (%q) does not exist", sourceFile)
				}
			} else if source != "consul" && source != "etcd" && source != "zookeeper" && source != "boltdb" {
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
			acmeFile := cmd.Flag("source-file").Value.String()

			switch source {
			case "file":
				config.BackendConfig = FileBackend{
					Name: FILE,
					Path: acmeFile,
				}
			case "consul":
				config.BackendConfig = KVBackend{
					Name:   CONSUL,
					Client: []string{"localhost:8500"},
					Config: &store.Config{},
				}
			case "etcd":
				config.BackendConfig = KVBackend{
					Name:   ETCD,
					Client: []string{"localhost:8500"},
					Config: &store.Config{},
				}
			case "zookeeper":
				config.BackendConfig = KVBackend{
					Name:   ZK,
					Client: []string{"localhost:8500"},
					Config: &store.Config{},
				}
			case "boltdb":
				config.BackendConfig = KVBackend{
					Name:   BOLTDB,
					Client: []string{"localhost:8500"},
					Config: &store.Config{},
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

	dumpCmd.Flags().String("source", "file", "Source type. One of 'file', 'consul', 'etcd', 'zookeeper', 'boltdb'.")
	dumpCmd.Flags().String("file", "./acme.json", "Path to 'acme.json' file if source type is 'file'")

	/* TODO implement this
	dumpCmd.Flags().String("kv.client")
	dumpCmd.Flags().String("kv.connection-timeout")
	dumpCmd.Flags().String("kv.sync-period")
	dumpCmd.Flags().String("kv.bucket")
	dumpCmd.Flags().Bool("kv.persist-connection")
	dumpCmd.Flags().String("kv.username")
	dumpCmd.Flags().String("kv.password")
	dumpCmd.Flags().String("kv.token")
	dumpCmd.Flags().String("kv.tls-cert-file")
	dumpCmd.Flags().String("kv.tls-key-file")
	dumpCmd.Flags().String("kv.tls-ca-cert-file")
	*/

	dumpCmd.Flags().Bool("watch", true, "Enable watching changes.")
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
