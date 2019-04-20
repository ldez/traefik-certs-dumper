// Copyright © 2019 ldez <lfernandez.dev@gmail.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// kvCmd represents the kv command
var kvCmd = &cobra.Command{
	Use:   "kv",
	Short: "TODO",
	Long:  `TODO`,
}

func init() {
	rootCmd.AddCommand(kvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
