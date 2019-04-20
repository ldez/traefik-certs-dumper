package cmd

import (
	"strconv"

	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/consul"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// consulCmd represents the consul command
var consulCmd = &cobra.Command{
	Use:   "consul",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  consulRun,
}

func init() {
	kvCmd.AddCommand(consulCmd)

	consulCmd.Flags().String("token", "", "Token for consul.")
}

func consulRun(cmd *cobra.Command, _ []string) error {
	// FIXME shared with file and all KVs
	dumpPath := cmd.Flag("dest").Value.String()

	crtInfo := dumper.FileInfo{
		Name: cmd.Flag("crt-name").Value.String(),
		Ext:  cmd.Flag("crt-ext").Value.String(),
	}

	keyInfo := dumper.FileInfo{
		Name: cmd.Flag("key-name").Value.String(),
		Ext:  cmd.Flag("key-ext").Value.String(),
	}

	subDir, _ := strconv.ParseBool(cmd.Flag("domain-subdir").Value.String())

	// ---

	config, err := getBaseConfig(cmd)
	if err != nil {
		return err
	}

	config.Options.Token = cmd.Flag("token").Value.String()

	config.Backend = store.CONSUL
	consul.Register()

	return kv.Dump(config, dumpPath, crtInfo, keyInfo, subDir)
}
