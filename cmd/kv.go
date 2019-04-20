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
	Short: "TODO",
	Long:  `TODO`,
}

func init() {
	rootCmd.AddCommand(kvCmd)

	kvCmd.PersistentFlags().StringSlice("endpoints", []string{"localhost:8500"}, "Comma separated list of endpoints.")
	kvCmd.PersistentFlags().Int("connection-timeout", 0, "Connection timeout in seconds.")
	kvCmd.PersistentFlags().String("password", "", "Password for connection.")
	kvCmd.PersistentFlags().String("username", "", "Username for connection.")

	// FIXME review TLS parts
	kvCmd.PersistentFlags().Bool("tls.enable", false, "Enable TLS encryption.")
	kvCmd.PersistentFlags().Bool("tls.insecureskipverify", false, "Trust unverified certificates if TLS is enabled.")
	kvCmd.PersistentFlags().String("tls.ca-cert-file", "", "Root CA file for certificate verification if TLS is enabled.")
}

func getBaseConfig(cmd *cobra.Command) (*kv.BaseConfig, error) {
	endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		return nil, err
	}

	connectionTimeout, err := cmd.Flags().GetInt("connection-timeout")
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, err
	}

	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, err
	}

	return &kv.BaseConfig{
		Endpoints: endpoints,
		Options: &store.Config{
			ClientTLS:         nil,
			TLS:               nil,
			ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
			SyncPeriod:        0,
			Bucket:            "",
			PersistConnection: false,
			Username:          username,
			Password:          password,
			Token:             "",
		},
	}, nil
}
