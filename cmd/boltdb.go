package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// boltdbCmd represents the boltdb command
var boltdbCmd = &cobra.Command{
	Use:   "boltdb",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("boltdb called")
	},
}

func init() {
	kvCmd.AddCommand(boltdbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// boltdbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boltdbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
