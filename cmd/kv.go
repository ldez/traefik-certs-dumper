package cmd

import (
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
