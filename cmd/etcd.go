package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// etcdCmd represents the etcd command
var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "TODO",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("etcd called")
		return nil
	},
}

func init() {
	kvCmd.AddCommand(etcdCmd)

	etcdCmd.Flags().Int("sync-period", 0, "Sync period for etcd in seconds.")
}
