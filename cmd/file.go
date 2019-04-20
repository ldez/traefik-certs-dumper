package cmd

import (
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/file"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: `Dump the content of the "acme.json" file.`,
	Long:  `Dump the content of the "acme.json" file from Traefik to certificates.`,
	RunE: runE(func(baseConfig *dumper.BaseConfig, cmd *cobra.Command) error {
		acmeFile := cmd.Flag("source").Value.String()

		return file.Dump(acmeFile, baseConfig)
	}),
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.Flags().String("source", "./acme.json", "Path to 'acme.json' file.")
}
