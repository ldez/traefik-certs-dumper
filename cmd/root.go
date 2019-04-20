package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ldez/traefik-certs-dumper/dumper"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "traefik-certs-dumper",
	Short: "Dump Let's Encrypt certificates from Traefik.",
	Long:  `Dump Let's Encrypt certificates from Traefik.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}

		crtExt := cmd.Flag("crt-ext").Value.String()
		keyExt := cmd.Flag("key-ext").Value.String()

		subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
		if !subDir {
			if crtExt == keyExt {
				return fmt.Errorf("--crt-ext (%q) and --key-ext (%q) are identical, in this case --domain-subdir is required", crtExt, keyExt)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.traefik-certs-dumper.yaml)")

	rootCmd.PersistentFlags().String("dest", "./dump", "Path to store the dump content.")
	rootCmd.PersistentFlags().String("crt-ext", ".crt", "The file extension of the generated certificates.")
	rootCmd.PersistentFlags().String("crt-name", "certificate", "The file name (without extension) of the generated certificates.")
	rootCmd.PersistentFlags().String("key-ext", ".key", "The file extension of the generated private keys.")
	rootCmd.PersistentFlags().String("key-name", "privatekey", "The file name (without extension) of the generated private keys.")
	rootCmd.PersistentFlags().Bool("domain-subdir", false, "Use domain as sub-directory.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".traefik-certs-dumper" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".traefik-certs-dumper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getBaseConfig(cmd *cobra.Command) (*dumper.BaseConfig, error) {
	subDir, err := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())
	if err != nil {
		return nil, err
	}

	return &dumper.BaseConfig{
		DumpPath: cmd.Flag("dest").Value.String(),
		CrtInfo: dumper.FileInfo{
			Name: cmd.Flag("crt-name").Value.String(),
			Ext:  cmd.Flag("crt-ext").Value.String(),
		},
		KeyInfo: dumper.FileInfo{
			Name: cmd.Flag("key-name").Value.String(),
			Ext:  cmd.Flag("key-ext").Value.String(),
		},
		DomainSubDir: subDir,
	}, nil
}

func runE(apply func(*dumper.BaseConfig, *cobra.Command) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := getBaseConfig(cmd)
		if err != nil {
			return err
		}

		err = apply(baseConfig, cmd)
		if err != nil {
			return err
		}

		return dumper.Tree(baseConfig.DumpPath, "")
	}
}
