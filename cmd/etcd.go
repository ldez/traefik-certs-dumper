package cmd

import (
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/etcd/v2"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// etcdCmd represents the etcd command
var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  runE(etcdRun),
}

func init() {
	kvCmd.AddCommand(etcdCmd)

	etcdCmd.Flags().Int("sync-period", 0, "Sync period for etcd in seconds.")
}

func etcdRun(baseConfig *dumper.BaseConfig, cmd *cobra.Command) error {
	config, err := getKvConfig(cmd)
	if err != nil {
		return err
	}

	synPeriod, err := cmd.Flags().GetInt("sync-period")
	if err != nil {
		return err
	}
	config.Options.SyncPeriod = time.Duration(synPeriod) * time.Second

	config.Backend = store.ETCD
	etcd.Register()

	return kv.Dump(config, baseConfig)
}
