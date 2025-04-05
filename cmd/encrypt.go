/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
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
		const encrypt = "encrypt"

		fileUtils := utils.NewFileUtils(args)
		pdfProcessor := pdf.NewPDFProcessor(fileUtils, encrypt)

		// check if any files/folders are provided
		pdfs, dir, interactive, err := fileUtils.CheckProvidedArgs()
		if handleError(cmd, err) {
			return
		}

		if interactive {
			selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, encrypt)
			if handleError(cmd, err) || quit {
				return
			}

			if len(selectedPdfs) == 0 {
				cmd.Println(infoStyle.Render("No PDFs were selected. Exiting."))
				return
			}
		}

		pword, quit, err = getPassword(cmd)
		if handleError(cmd, err) || quit {
			return
		}

		if !interactive {
			selectedPdfs = pdfs
		}

		saveDir, err := fileUtils.GetCurrentWorkingDir()
		if handleError(cmd, err) {
			return
		}

		processPDFs(cmd, pdfProcessor, selectedPdfs, dir, saveDir, pword)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF files.")
	if err := encryptCmd.Parent().MarkFlagRequired("password"); err != nil {
		return
	}

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
