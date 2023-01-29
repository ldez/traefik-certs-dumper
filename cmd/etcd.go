package cmd

import (
	"context"
	"time"

	"github.com/kvtools/etcdv2"
	"github.com/kvtools/etcdv3"
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

	backend, err := cmd.Flags().GetString("etcd-version")
	if err != nil {
		return err
	}

	tlsConfig, err := createTLSConfig(cmd)
	if err != nil {
		return err
	}

	synPeriod, err := cmd.Flags().GetInt("sync-period")
	if err != nil {
		return err
	}

	connectionTimeout, err := cmd.Flags().GetInt("connection-timeout")
	if err != nil {
		return err
	}

	switch backend {
	case "etcdv3":
		config.Options = &etcdv3.Config{
			TLS:               tlsConfig,
			ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
			SyncPeriod:        time.Duration(synPeriod) * time.Second,
			Username:          cmd.Flag("password").Value.String(),
			Password:          cmd.Flag("username").Value.String(),
		}
		config.StoreName = etcdv3.StoreName
	default:
		config.Options = &etcdv2.Config{
			TLS:               tlsConfig,
			ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
			SyncPeriod:        time.Duration(synPeriod) * time.Second,
			Username:          cmd.Flag("password").Value.String(),
			Password:          cmd.Flag("username").Value.String(),
		}

		config.StoreName = etcdv2.StoreName
	}

	return kv.Dump(context.Background(), config, baseConfig)
}
