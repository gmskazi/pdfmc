/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [files... or folder]",
	Short: "Encrypt PDF files.",
	Long:  `This is a tool for encrypting PDF files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			selectedPdfs []string
			pword        string
			quit         bool
		)

		pword, err := cmd.Flags().GetString("password")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		fileUtils := utils.NewFileUtils(args)
		pdfProcessor := pdf.NewPDFProcessor(fileUtils, "encrypt")

		// check if any files/folders are provided
		pdfs, dir, interactive, err := fileUtils.CheckProvidedArgs()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		if interactive {
			selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, "encrypt")
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			if len(selectedPdfs) == 0 {
				cmd.PrintErrln(infoStyle.Render("No PDFs were selected. Exiting."))
				return
			}
		}

		if pword == "" {
			pword, quit, err = textInputs.TextinputInteractive()
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			fmt.Println()
		}

		if !interactive {
			selectedPdfs = pdfs
		}

		saveDir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		for _, pdf := range selectedPdfs {
			encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword)
			if err != nil {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			complete := fmt.Sprintf("PDF file encrypted successfully to: %s/%s", saveDir, encryptedPdf)
			cmd.Println(selectedStyle.Render(complete))
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF files.")

	// autocomplete for files flag
	mergeCmd.ValidArgsFunction = GetSuggestions
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
