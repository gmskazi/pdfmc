/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"github.com/gmskazi/pdfmc/cmd/autocomplete"
	"github.com/gmskazi/pdfmc/cmd/program"
	"github.com/gmskazi/pdfmc/cmd/styles"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [files... or folder]",
	Short: "Encrypt PDF files.",
	Long:  `This is a tool for encrypting PDF files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := program.NewProgram(cmd, args, encrypt)
		if err := p.ExecuteEncrypt(); err != nil {
			cmd.PrintErrln(styles.ErrorStyle.Render(err.Error()))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF files.")

	// autocomplete for files flag
	mergeCmd.ValidArgsFunction = autocomplete.GetSuggestions
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
