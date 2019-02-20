package main

import (
	"fmt"
	"log"
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
			crtExt := cmd.Flag("crt-ext").Value.String()
			keyExt := cmd.Flag("key-ext").Value.String()
			subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
			if crtExt == keyExt && !subDir {
				return fmt.Errorf("--crt-ext (%q) and --key-ext (%q) are identical, in this case --domain-subdir is required", crtExt, keyExt)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			acmeFile := cmd.Flag("source").Value.String()
			dumpPath := cmd.Flag("dest").Value.String()
			crtExt := cmd.Flag("crt-ext").Value.String()
			keyExt := cmd.Flag("key-ext").Value.String()
			subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())

			err := dump(acmeFile, dumpPath, crtExt, keyExt, subDir)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	dumpCmd.Flags().String("source", "./acme.json", "Path to 'acme.json' file.")
	dumpCmd.Flags().String("dest", "./dump", "Path to store the dump content.")
	dumpCmd.Flags().String("crt-ext", ".crt", "The file extension of the generated certificates.")
	dumpCmd.Flags().String("key-ext", ".key", "The file extension of the generated private keys.")
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
