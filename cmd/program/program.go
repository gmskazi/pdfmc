package program

import (
	"errors"
	"fmt"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/styles"
	"github.com/gmskazi/pdfmc/cmd/ui/multiReorder"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Program struct {
	cmd   *cobra.Command
	args  []string
	logo  string
	name  string
	pword string
	MergeFlags
}

type MergeFlags struct {
	reorder bool
	encrypt bool
}

func NewProgram(cmd *cobra.Command, args []string, logo string) *Program {
	mergeFlags := MergeFlags{
		reorder: getFlagBoolValue(cmd, "order"),
		encrypt: getFlagBoolValue(cmd, "encrypt"),
	}

	return &Program{
		cmd:        cmd,
		args:       args,
		logo:       logo,
		name:       getFlagValue(cmd.Flag("name")),
		pword:      getFlagValue(cmd.Flag("password")),
		MergeFlags: mergeFlags,
	}
}

func getFlagValue(flag *pflag.Flag) string {
	if flag == nil {
		return ""
	}
	return flag.Value.String()
}

func getFlagBoolValue(cmd *cobra.Command, flagname string) bool {
	value, err := cmd.Flags().GetBool(flagname)
	if err != nil {
		return false
	}
	return value
}

func (p *Program) getPassword() error {
	// check and update the password
	if p.pword == "" {
		newPword, quit, err := textInputs.TextinputInteractive()
		if err != nil || quit {
			return err
		}
		p.pword = newPword
		return nil
	}
	return nil
}

func (p *Program) processEncryptPDFs(pdfProcessor *pdf.PDFProcessor, selectedPdfs []string, dir, saveDir, pword string) error {
	for _, pdf := range selectedPdfs {
		if p.logo == "merge" {
			encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword, "")
			if err != nil {
				return err
			}
			complete := fmt.Sprintf("PDF files merged and encrypted successfully to: %s/%s", saveDir, encryptedPdf)
			p.cmd.Println(styles.SelectedStyle.Render(complete))
		} else {
			encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword, p.name)
			if err != nil {
				return err
			}
			complete := fmt.Sprintf("PDF file encrypted successfully to: %s/%s", saveDir, encryptedPdf)
			p.cmd.Println(styles.SelectedStyle.Render(complete))
		}
	}
	return nil
}

func (p *Program) processDecryptPDFs(pdfProcessor *pdf.PDFProcessor, selectedPdfs []string, dir, saveDir, pword string) error {
	for _, pdf := range selectedPdfs {
		encryptedPdf, err := pdfProcessor.DecryptPdf(pdf, dir, pword, p.name)
		if err != nil {
			return err
		}

		complete := fmt.Sprintf("PDF file decrypted successfully to: %s/%s", saveDir, encryptedPdf)
		p.cmd.Println(styles.SelectedStyle.Render(complete))
	}
	return nil
}

func (p *Program) ExecuteEncrypt() error {
	var (
		selectedPdfs []string
		quit         bool
		err          error
	)

	f := utils.NewFileUtils(p.args)
	pdfProcessor := pdf.NewPDFProcessor(p.logo)

	pdfs, dir, err := f.CheckProvidedArgs()
	if err != nil {
		return err
	}

	if f.Interactive {
		selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, p.logo)
		if err != nil || quit {
			return err
		}

		if len(selectedPdfs) == 0 {
			p.cmd.Println(styles.InfoStyle.Render("No PDFs were selected. Exiting."))
			return err
		}
	} else {
		selectedPdfs = pdfs
	}

	if err := p.getPassword(); err != nil {
		return err
	}

	saveDir, err := f.GetCurrentWorkingDir()
	if err != nil {
		return err
	}

	if err := p.processEncryptPDFs(pdfProcessor, selectedPdfs, dir, saveDir, p.pword); err != nil {
		return err
	}
	return nil
}

func (p *Program) ExecuteMerge() error {
	var (
		selectedPdfs []string
		quit         bool
		err          error
	)

	if p.encrypt && p.pword != "" {
		return errors.New("please provide either the --password flag or use the --encrypt flag for interactive encryption")
	}
	f := utils.NewFileUtils(p.args)

	// check if any files/folders are provided
	pdfs, dir, err := f.CheckProvidedArgs()
	if err != nil {
		return err
	}

	if f.Interactive {
		for {
			selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, p.logo)
			if err != nil || quit {
				return err
			}

			if len(selectedPdfs) <= 1 {
				continue
			}
			break
		}
	} else {
		selectedPdfs = pdfs
	}

	// reordering of the pdfs
	if p.reorder {
		selectedPdfs, quit, err = multiReorder.MultiReorderInteractive(selectedPdfs, p.logo)
		if err != nil || quit {
			return err
		}
	}

	pdfWithFullPath := f.AddFullPathToPdfs(dir, selectedPdfs)

	pdfProcessor := pdf.NewPDFProcessor(p.logo)

	p.name, err = pdfProcessor.MergePdfs(pdfWithFullPath, p.name)
	if err != nil {
		return err
	}

	saveDir, err := f.GetCurrentWorkingDir()
	if err != nil {
		return err
	}

	// if the encrypt flag is set, ask for password interactively
	if p.encrypt {
		p.pword, quit, err = textInputs.TextinputInteractive()
		if err != nil || quit {
			return err
		}

		fmt.Println()
	}

	// encrypt pdf file if flag is set
	if p.pword != "" {
		if err := p.processEncryptPDFs(pdfProcessor, []string{p.name}, saveDir, saveDir, p.pword); err != nil {
			return err
		}
	} else {
		complete := fmt.Sprintf("PDF files merged successfully to: %s/%s", saveDir, p.name)
		p.cmd.Println(styles.InfoStyle.Render(complete))
	}

	return nil
}

func (p *Program) ExecuteDecrypt() error {
	var (
		selectedPdfs []string
		quit         bool
		err          error
	)

	f := utils.NewFileUtils(p.args)
	pdfProcessor := pdf.NewPDFProcessor(p.logo)

	pdfs, dir, err := f.CheckProvidedArgs()
	if err != nil {
		return err
	}

	if f.Interactive {
		selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, p.logo)
		if err != nil || quit {
			return err
		}

		if len(selectedPdfs) == 0 {
			p.cmd.Println(styles.InfoStyle.Render("No PDFs were selected. Exiting."))
			return err
		}
	} else {
		selectedPdfs = pdfs
	}

	if err := p.getPassword(); err != nil {
		return err
	}

	saveDir, err := f.GetCurrentWorkingDir()
	if err != nil {
		return err
	}

	if err := p.processDecryptPDFs(pdfProcessor, selectedPdfs, dir, saveDir, p.pword); err != nil {
		return err
	}
	return nil
}
