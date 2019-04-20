package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// zookeeperCmd represents the zookeeper command
var zookeeperCmd = &cobra.Command{
	Use:   "zookeeper",
	Short: "TODO",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("zookeeper called")
		return nil
	},
}

func init() {
	kvCmd.AddCommand(zookeeperCmd)
}
