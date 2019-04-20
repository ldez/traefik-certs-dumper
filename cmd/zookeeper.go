package cmd

import (
	"strconv"

	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/zookeeper"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// zookeeperCmd represents the zookeeper command
var zookeeperCmd = &cobra.Command{
	Use:   "zookeeper",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  zookeeperRun,
}

func init() {
	kvCmd.AddCommand(zookeeperCmd)
}

func zookeeperRun(cmd *cobra.Command, _ []string) error {
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

	config.Backend = store.ZK
	zookeeper.Register()

	return kv.Dump(config, dumpPath, crtInfo, keyInfo, subDir)
}
