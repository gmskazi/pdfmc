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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt [files... or folder]",
	Short: "Decrypt PDF files.",
	Long:  `This is a tool to decrypt pdf files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := program.NewProgram(cmd, args, decrypt)
		if err := p.ExecuteDecrypt(); err != nil {
			cmd.PrintErrln(styles.ErrorStyle.Render(err.Error()))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("password", "p", "", "Password to decrypt the PDF files.")
	decryptCmd.Flags().StringP("name", "n", "", "Add a prefix to the beginning of the file name.")
	// autocomplete for files
	decryptCmd.ValidArgsFunction = autocomplete.GetSuggestions

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
