package cmd

import (
	"strconv"

	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/ldez/traefik-certs-dumper/dumper"
	"github.com/ldez/traefik-certs-dumper/dumper/kv"
	"github.com/spf13/cobra"
)

// boltdbCmd represents the boltdb command
var boltdbCmd = &cobra.Command{
	Use:   "boltdb",
	Short: "TODO",
	Long:  `TODO`,
	RunE:  boltdbRun,
}

func init() {
	kvCmd.AddCommand(boltdbCmd)

	boltdbCmd.Flags().Bool("persist-connection", false, "Persist connection for boltdb.")
	boltdbCmd.Flags().String("bucket", "traefik", "Bucket for boltdb.")
}

func boltdbRun(cmd *cobra.Command, _ []string) error {
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

	config.Options.Bucket = cmd.Flag("bucket").Value.String()
	config.Options.PersistConnection, _ = cmd.Flags().GetBool("persist-connection")

	config.Backend = store.BOLTDB
	boltdb.Register()

	return kv.Dump(config, dumpPath, crtInfo, keyInfo, subDir)
}
