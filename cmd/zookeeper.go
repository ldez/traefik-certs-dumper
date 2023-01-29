package cmd

import (
	"context"
	"time"

	"github.com/kvtools/zookeeper"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/ldez/traefik-certs-dumper/v2/dumper/kv"
	"github.com/spf13/cobra"
)

// zookeeperCmd represents the zookeeper command.
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

	connectionTimeout, err := cmd.Flags().GetInt("connection-timeout")
	if err != nil {
		return err
	}

	config.Options = &zookeeper.Config{
		ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
		Username:          cmd.Flag("password").Value.String(),
		Password:          cmd.Flag("username").Value.String(),
		MaxBufferSize:     0,
	}

	config.StoreName = zookeeper.StoreName

	return kv.Dump(context.Background(), config, baseConfig)
}
