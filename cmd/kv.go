package cmd

import (
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// kvCmd represents the kv command
var kvCmd = &cobra.Command{
	Use:   "kv",
	Short: `Dump the content of a KV store.`,
	Long:  `Dump the content of a KV store.`,
}

func init() {
	rootCmd.AddCommand(kvCmd)

	kvCmd.PersistentFlags().StringSlice("endpoints", []string{"localhost:8500"}, "Comma separated list of endpoints.")
	kvCmd.PersistentFlags().Int("connection-timeout", 0, "Connection timeout in seconds.")
	kvCmd.PersistentFlags().String("prefix", "traefik", "Prefix used for KV store.")
	kvCmd.PersistentFlags().String("password", "", "Password for connection.")
	kvCmd.PersistentFlags().String("username", "", "Username for connection.")
	kvCmd.PersistentFlags().Bool("watch", false, "Enable watching changes.")

	// FIXME review TLS parts
	// kvCmd.PersistentFlags().Bool("tls.enable", false, "Enable TLS encryption.")
	// kvCmd.PersistentFlags().Bool("tls.insecureskipverify", false, "Trust unverified certificates if TLS is enabled.")
	// kvCmd.PersistentFlags().String("tls.ca-cert-file", "", "Root CA file for certificate verification if TLS is enabled.")
}

func getKvConfig(cmd *cobra.Command) (*kv.Config, error) {
	endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		return nil, err
	}

	connectionTimeout, err := cmd.Flags().GetInt("connection-timeout")
	if err != nil {
		return nil, err
	}

	return &kv.Config{
		Endpoints: endpoints,
		Prefix:    cmd.Flag("prefix").Value.String(),
		Options: &store.Config{
			ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
			Username:          cmd.Flag("password").Value.String(),
			Password:          cmd.Flag("username").Value.String(),
		},
	}, nil
}
