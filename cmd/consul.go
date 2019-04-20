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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("consul called")
		return nil
	},
}

func init() {
	kvCmd.AddCommand(consulCmd)

	consulCmd.Flags().String("token", "", "Token for consul.")
}
