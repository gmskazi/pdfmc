package cmd

import (
	"fmt"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/spf13/cobra"
)

func handleError(cmd *cobra.Command, err error) error {
	if err != nil {
		cmd.PrintErrln(errorStyle.Render(err.Error()))
		return err
	}
	return nil
}

func getPassword(cmd *cobra.Command) (string, error) {
	flagPassword, err := cmd.Flags().GetString("password")
	if err != nil {
		return "", err
	}
	if flagPassword == "" {
		pword, quit, err := textInputs.TextinputInteractive()
		if err != nil || quit {
			return "", err
		}
		return pword, nil
	}
	return flagPassword, nil
}

func processPDFs(cmd *cobra.Command, pdfProcessor *pdf.PDFProcessor, selectedPdfs []string, dir, saveDir, pword string) {
	for _, pdf := range selectedPdfs {
		encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword)
		if err = handleError(cmd, err); err != nil {
			return
		}

		complete := fmt.Sprintf("PDF file encrypted successfully to: %s/%s", saveDir, encryptedPdf)
		cmd.Println(selectedStyle.Render(complete))
	}
}

func executeEncrypt(cmd *cobra.Command, args []string) error {
	var (
		selectedPdfs []string
		pword        string
		quit         bool
	)

	fileUtils := utils.NewFileUtils(args)
	pdfProcessor := pdf.NewPDFProcessor(fileUtils, "encrypt")

	pdfs, dir, interactive, err := fileUtils.CheckProvidedArgs()
	if err = handleError(cmd, err); err != nil {
		return err
	}

	if interactive {
		selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, "encrypt")
		if err = handleError(cmd, err); err != nil || quit {
			return err
		}

		if len(selectedPdfs) == 0 {
			cmd.Println(infoStyle.Render("No PDFs were selected. Exiting."))
			return err
		}
	}

	pword, err = getPassword(cmd)
	if err = handleError(cmd, err); err != nil {
		return err
	}

	if !interactive {
		selectedPdfs = pdfs
	}

	saveDir, err := fileUtils.GetCurrentWorkingDir()
	if err = handleError(cmd, err); err != nil {
		return err
	}

	processPDFs(cmd, pdfProcessor, selectedPdfs, dir, saveDir, pword)
	return nil
}
