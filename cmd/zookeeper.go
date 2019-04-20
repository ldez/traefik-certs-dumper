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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper called")
	},
}

func init() {
	kvCmd.AddCommand(zookeeperCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// zookeeperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// zookeeperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
