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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("etcd called")
	},
}

func init() {
	kvCmd.AddCommand(etcdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// etcdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// etcdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
