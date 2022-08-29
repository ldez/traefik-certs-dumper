package cmd

import (
	"context"

	"github.com/kvtools/valkeyrie/store"
	"github.com/kvtools/valkeyrie/store/consul"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/ldez/traefik-certs-dumper/v2/dumper/kv"
	"github.com/spf13/cobra"
)

// consulCmd represents the consul command.
var consulCmd = &cobra.Command{
	Use:   "consul",
	Short: "Dump the content of Consul.",
	Long:  `Dump the content of Consul.`,
	RunE:  runE(consulRun),
}

func init() {
	kvCmd.AddCommand(consulCmd)

	consulCmd.Flags().String("token", "", "Token for consul.")
}

func consulRun(baseConfig *dumper.BaseConfig, cmd *cobra.Command) error {
	config, err := getKvConfig(cmd)
	if err != nil {
		return err
	}

	config.Options.Token = cmd.Flag("token").Value.String()

	config.Backend = store.CONSUL
	consul.Register()

	return kv.Dump(context.Background(), config, baseConfig)
}
