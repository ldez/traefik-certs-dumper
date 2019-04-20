package cmd

import (
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// boltdbCmd represents the boltdb command
var boltdbCmd = &cobra.Command{
	Use:   "boltdb",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  runE(boltdbRun),
}

func init() {
	kvCmd.AddCommand(boltdbCmd)

	boltdbCmd.Flags().Bool("persist-connection", false, "Persist connection for boltdb.")
	boltdbCmd.Flags().String("bucket", "traefik", "Bucket for boltdb.")
}

func boltdbRun(baseConfig *dumper.BaseConfig, cmd *cobra.Command) error {
	config, err := getKvConfig(cmd)
	if err != nil {
		return err
	}

	config.Options.Bucket = cmd.Flag("bucket").Value.String()
	config.Options.PersistConnection, _ = cmd.Flags().GetBool("persist-connection")

	config.Backend = store.BOLTDB
	boltdb.Register()

	return kv.Dump(config, baseConfig)
}
