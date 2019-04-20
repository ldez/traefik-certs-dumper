package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// consulCmd represents the consul command
var consulCmd = &cobra.Command{
	Use:   "consul",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consul called")
	},
}

func init() {
	kvCmd.AddCommand(consulCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consulCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consulCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
