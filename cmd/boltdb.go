package cmd

import (
	"context"
	"time"

	"github.com/kvtools/boltdb"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/ldez/traefik-certs-dumper/v2/dumper/kv"
	"github.com/spf13/cobra"
)

// boltdbCmd represents the boltdb command.
var boltdbCmd = &cobra.Command{
	Use:   "boltdb",
	Short: "Dump the content of BoltDB.",
	Long:  `Dump the content of BoltDB.`,
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

	connectionTimeout, err := cmd.Flags().GetInt("connection-timeout")
	if err != nil {
		return err
	}

	persistConnection, _ := cmd.Flags().GetBool("persist-connection")

	config.Options = &boltdb.Config{
		Bucket:            cmd.Flag("bucket").Value.String(),
		PersistConnection: persistConnection,
		ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
	}

	config.StoreName = boltdb.StoreName

	return kv.Dump(context.Background(), config, baseConfig)
}
