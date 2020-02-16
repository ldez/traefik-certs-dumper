package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/ldez/traefik-certs-dumper/v2/dumper/kv"
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

	kvCmd.PersistentFlags().StringSlice("endpoints", []string{"localhost:8500"}, "List of endpoints.")
	kvCmd.PersistentFlags().Int("connection-timeout", 0, "Connection timeout in seconds.")
	kvCmd.PersistentFlags().String("prefix", "traefik", "Prefix used for KV store.")
	kvCmd.PersistentFlags().String("suffix", kv.DefaultStoreKeySuffix, "Suffix/Storage used for KV store.")
	kvCmd.PersistentFlags().String("password", "", "Password for connection.")
	kvCmd.PersistentFlags().String("username", "", "Username for connection.")

	kvCmd.PersistentFlags().Bool("tls", false, "Enable TLS encryption.")
	kvCmd.PersistentFlags().String("tls.ca", "", "Root CA for certificate verification if TLS is enabled")
	kvCmd.PersistentFlags().Bool("tls.ca.optional", false, "")
	kvCmd.PersistentFlags().String("tls.cert", "", "TLS cert")
	kvCmd.PersistentFlags().String("tls.key", "", "TLS key")
	kvCmd.PersistentFlags().Bool("tls.insecureskipverify", false, "Trust unverified certificates if TLS is enabled.")
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

	tlsConfig, err := createTLSConfig(cmd)
	if err != nil {
		return nil, err
	}

	return &kv.Config{
		Endpoints: endpoints,
		Prefix:    cmd.Flag("prefix").Value.String(),
		Suffix:    cmd.Flag("suffix").Value.String(),
		Options: &store.Config{
			ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
			Username:          cmd.Flag("password").Value.String(),
			Password:          cmd.Flag("username").Value.String(),
			TLS:               tlsConfig,
		},
	}, nil
}

func createTLSConfig(cmd *cobra.Command) (*tls.Config, error) {
	enable, _ := cmd.Flags().GetBool("tls")
	if !enable {
		return nil, nil
	}

	ca := cmd.Flag("tls.ca").Value.String()
	caPool, err := getCertPool(ca)
	if err != nil {
		return nil, err
	}

	caOptional, _ := cmd.Flags().GetBool("tls.ca.optional")
	clientAuth := getClientAuth(ca, caOptional)

	insecureSkipVerify, _ := cmd.Flags().GetBool("tls.insecureskipverify")
	privateKey := cmd.Flag("tls.key").Value.String()
	certContent := cmd.Flag("tls.cert").Value.String()

	if !insecureSkipVerify && (len(certContent) == 0 || len(privateKey) == 0) {
		return nil, fmt.Errorf("TLS Certificate or Key file must be set when TLS configuration is created")
	}

	cert, err := getCertificate(privateKey, certContent)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS keypair: %w", err)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: insecureSkipVerify,
		ClientAuth:         clientAuth,
	}, nil
}

func getCertPool(ca string) (*x509.CertPool, error) {
	caPool := x509.NewCertPool()

	if ca != "" {
		caContent, err := getCAContent(ca)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA. %w", err)
		}

		if !caPool.AppendCertsFromPEM(caContent) {
			return nil, fmt.Errorf("failed to parse CA")
		}
	}

	return caPool, nil
}

func getCAContent(ca string) ([]byte, error) {
	if _, errCA := os.Stat(ca); errCA != nil {
		return []byte(ca), nil
	}

	caContent, err := ioutil.ReadFile(filepath.Clean(ca))
	if err != nil {
		return nil, err
	}
	return caContent, nil
}

func getClientAuth(ca string, caOptional bool) tls.ClientAuthType {
	if ca == "" {
		return tls.NoClientCert
	}

	if caOptional {
		return tls.VerifyClientCertIfGiven
	}
	return tls.RequireAndVerifyClientCert
}

func getCertificate(privateKey, certContent string) (tls.Certificate, error) {
	if certContent == "" || privateKey == "" {
		return tls.Certificate{}, nil
	}

	_, errKeyIsFile := os.Stat(privateKey)
	_, errCertIsFile := os.Stat(certContent)

	if errCertIsFile == nil && os.IsNotExist(errKeyIsFile) {
		return tls.Certificate{}, fmt.Errorf("tls cert is a file, but tls key is not")
	}

	if os.IsNotExist(errCertIsFile) && errKeyIsFile == nil {
		return tls.Certificate{}, fmt.Errorf("TLS key is a file, but tls cert is not")
	}

	// string
	if os.IsNotExist(errCertIsFile) && os.IsNotExist(errKeyIsFile) {
		return tls.X509KeyPair([]byte(certContent), []byte(privateKey))
	}

	// files
	if errCertIsFile == nil && errKeyIsFile == nil {
		return tls.LoadX509KeyPair(certContent, privateKey)
	}

	if errCertIsFile != nil {
		return tls.Certificate{}, errCertIsFile
	}
	return tls.Certificate{}, errKeyIsFile
}
