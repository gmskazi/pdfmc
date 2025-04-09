/*
Copyright © 2025 Aito Nakajima
*/
package cmd

import (
	"github.com/gmskazi/pdfmc/cmd/autocomplete"
	"github.com/gmskazi/pdfmc/cmd/program"
	"github.com/gmskazi/pdfmc/cmd/styles"
	"github.com/spf13/cobra"
)

const (
	merge   = "merge"
	encrypt = "encrypt"
	decrypt = "decrypt"
)

var name string

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [files... or folder]",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := program.NewProgram(cmd, args, merge)
		if err := p.ExecuteMerge(); err != nil {
			cmd.PrintErrln(styles.ErrorStyle.Render(err.Error()))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&name, "name", "n", "merged_output", "Custom name for the merged PDF files")
	mergeCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF file.")
	mergeCmd.Flags().BoolP("order", "o", false, "Reorder the PDF files before merging.")
	mergeCmd.Flags().BoolP("encrypt", "e", false, "Encrypt the PDF file interatively.")

	// autocomplete for files flag
	mergeCmd.ValidArgsFunction = autocomplete.GetSuggestions
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
