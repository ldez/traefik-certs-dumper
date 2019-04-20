package cmd

import (
	"strconv"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/etcd/v2"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// etcdCmd represents the etcd command
var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  etcdRun,
}

func init() {
	kvCmd.AddCommand(etcdCmd)

	etcdCmd.Flags().Int("sync-period", 0, "Sync period for etcd in seconds.")
}

func etcdRun(cmd *cobra.Command, _ []string) error {
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

	synPeriod, err := cmd.Flags().GetInt("sync-period")
	if err != nil {
		return err
	}
	config.Options.SyncPeriod = time.Duration(synPeriod) * time.Second

	config.Backend = store.ETCD
	etcd.Register()

	return kv.Dump(config, dumpPath, crtInfo, keyInfo, subDir)
}
