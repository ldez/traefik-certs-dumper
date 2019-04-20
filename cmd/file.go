package cmd

import (
	"strconv"

	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/file"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: `Dump the content of the "acme.json" file.`,
	Long:  `Dump the content of the "acme.json" file from Traefik to certificates.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		acmeFile := cmd.Flag("source").Value.String()
		dumpPath := cmd.Flag("dest").Value.String()

		crtInfo := dumper.FileInfo{
			Name: cmd.Flag("crt-name").Value.String(),
			Ext:  cmd.Flag("crt-ext").Value.String(),
		}

		keyInfo := dumper.FileInfo{
			Name: cmd.Flag("key-name").Value.String(),
			Ext:  cmd.Flag("key-ext").Value.String(),
		}

		subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())

		err := file.Dump(acmeFile, dumpPath, crtInfo, keyInfo, subDir)
		if err != nil {
			return err
		}

		return dumper.Tree(dumpPath, "")
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.Flags().String("source", "./acme.json", "Path to 'acme.json' file.")
}
