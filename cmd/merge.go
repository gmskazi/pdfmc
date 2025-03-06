/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"github.com/gmskazi/pdfmergecrypt/cmd/utils"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Hello()
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
