package cmd

import (
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/zookeeper"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// zookeeperCmd represents the zookeeper command
var zookeeperCmd = &cobra.Command{
	Use:   "zookeeper",
	Short: "Dump the content of zookeeper.",
	Long:  `Dump the content of zookeeper.`,
	RunE:  runE(zookeeperRun),
}

func init() {
	kvCmd.AddCommand(zookeeperCmd)
}

func zookeeperRun(baseConfig *dumper.BaseConfig, cmd *cobra.Command) error {
	config, err := getKvConfig(cmd)
	if err != nil {
		return err
	}

	config.Backend = store.ZK
	zookeeper.Register()

	return kv.Dump(config, baseConfig)
}
