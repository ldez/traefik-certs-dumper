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

	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump Let's Encrypt certificates from Traefik",
		Long:  `Dump the content of the "acme.json" file from Traefik to certificates.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			source := cmd.Flag("source").Value.String()
			if source != "file" && source != "consul" {
				return fmt.Errorf("--source (%q) is not allowed, use one of 'file' or 'consul'", source)
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
			acmeFile := cmd.Flag("source-file").Value.String()
			dumpPath := cmd.Flag("dest").Value.String()

			crtInfo := fileInfo{
				Name: cmd.Flag("crt-name").Value.String(),
				Ext:  cmd.Flag("crt-ext").Value.String(),
			}

			keyInfo := fileInfo{
				Name: cmd.Flag("key-name").Value.String(),
				Ext:  cmd.Flag("key-ext").Value.String(),
			}

			subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
			watchConsul, _ := strconv.ParseBool(cmd.Flag("source-consul-watch").Value.String())

			switch cmd.Flag("source").Value.String() {

			case "consul":
				dumpConsul(watchConsul, dumpPath, crtInfo, keyInfo, subDir)

			case "file":
			default:
				err := dumpFile(acmeFile, dumpPath, crtInfo, keyInfo, subDir)
				if err != nil {
					return err
				}
				return tree(dumpPath, "")
			}
			return nil
		},
	}

	dumpCmd.Flags().String("source", "file", "Source type. One of 'file' or 'consul'. Consul connection parameters can be set via environment variables, see https://www.consul.io/docs/commands/index.html#environment-variables")
	dumpCmd.Flags().String("source-file", "./acme.json", "Path to 'acme.json' file if source type is 'file'")
	dumpCmd.Flags().Bool("source-consul-watch", true, "Enable watching changes in Consul.")
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
