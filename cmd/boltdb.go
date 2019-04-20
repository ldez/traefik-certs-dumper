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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("boltdb called")
		return nil
	},
}

func init() {
	kvCmd.AddCommand(boltdbCmd)

	boltdbCmd.Flags().Bool("persist-connection", false, "Persist connection for boltdb.")
	boltdbCmd.Flags().String("bucket", "traefik", "Bucket for boltdb.")
}
