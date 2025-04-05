package cmd

import (
	"fmt"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/spf13/cobra"
)

func handleError(cmd *cobra.Command, err error) bool {
	if err != nil {
		cmd.PrintErrln(errorStyle.Render(err.Error()))
		return true
	}
	return false
}

func getPassword(cmd *cobra.Command) (string, bool, error) {
	flagPassword, err := cmd.Flags().GetString("password")
	if err != nil {
		return "", false, err
	}
	if flagPassword == "" {
		pword, quit, err := textInputs.TextinputInteractive()
		if err != nil {
			return "", false, err
		}
		return pword, quit, nil
	}
	return flagPassword, false, nil
}

func processPDFs(cmd *cobra.Command, pdfProcessor *pdf.PDFProcessor, selectedPdfs []string, dir, saveDir, pword string) {
	for _, pdf := range selectedPdfs {
		encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword)
		if handleError(cmd, err) {
			return
		}

		complete := fmt.Sprintf("PDF file encrypted successfully to: %s/%s", saveDir, encryptedPdf)
		cmd.Println(selectedStyle.Render(complete))
	}
}
