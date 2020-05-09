package cmd

import (
	"time"

	"github.com/abronan/valkeyrie/store"
	etcdv2 "github.com/abronan/valkeyrie/store/etcd/v2"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/ldez/traefik-certs-dumper/v2/dumper/kv"
	"github.com/spf13/cobra"
)

// etcdCmd represents the etcd command.
var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "Dump the content of etcd.",
	Long:  `Dump the content of etcd.`,
	RunE:  runE(etcdRun),
}

func init() {
	kvCmd.AddCommand(etcdCmd)

	etcdCmd.Flags().Int("sync-period", 0, "Sync period for etcd in seconds.")
	etcdCmd.Flags().String("etcd-version", "etcd", "The etcd version can be: 'etcd' or 'etcdv3'.")
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

	backend, err := cmd.Flags().GetString("etcd-version")
	if err != nil {
		return err
	}

	switch backend {
	case "etcdv3":
		config.Backend = store.ETCDV3
		etcdv3.Register()
	default:
		config.Backend = store.ETCD
		etcdv2.Register()
	}

	return kv.Dump(config, baseConfig)
}
