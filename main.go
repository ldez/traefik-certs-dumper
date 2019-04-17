package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "traefik-certs-dumper",
		Short:   "Dump Let's Encrypt certificates from Traefik",
		Long:    `Dump the content of the "acme.json" file from Traefik to certificates.`,
		Version: version,
	}

	dumpConfig := &dumpConfig{}

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

			var backend Backend
			switch source {
			case "file":
				backend = FILE
			case "consul":
				backend = CONSUL
			case "etcd":
				backend = ETCD
			case "zookeeper":
				backend = ZK
			case "boltdb":
				backend = BOLTDB
			}

			dumpConfig.Path = cmd.Flag("dest").Value.String()

			dumpConfig.CertInfo = fileInfo{
				Name: cmd.Flag("crt-name").Value.String(),
				Ext:  cmd.Flag("crt-ext").Value.String(),
			}

			dumpConfig.KeyInfo = fileInfo{
				Name: cmd.Flag("key-name").Value.String(),
				Ext:  cmd.Flag("key-ext").Value.String(),
			}

			dumpConfig.DomainSubDir, _ = strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
			dumpConfig.Watch, _ = strconv.ParseBool(cmd.Flag("watch").Value.String())

			fmt.Println(dumpConfig)

			if backend == FILE {
				data, err := getAcmeDataFromJSONFile(acmeFile)
				if err != nil {
					return fmt.Errorf("[ERR] %v", err)
				}
				if err := dump(dumpConfig, data); err != nil {
					return err
				}
			} else if err := loop(dumpConfig, backend); err != nil {
				return err
			}

			return nil
		},
	}

	// TODO fill readme
	dumpCmd.Flags().String("source", "file", "Source type. One of 'file', 'consul', 'etcd', 'zookeeper', 'boltdb'. For configuration options of the Key/Value stores see https://github.com/ldez/traefik-certs-dumper#configuration-of-key-value-stores.")
	dumpCmd.Flags().String("source-file", "./acme.json", "Path to 'acme.json' file if source type is 'file'")
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
