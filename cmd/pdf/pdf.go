package pdf

import (
	"errors"

	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFProcessor struct {
	FileUtils *utils.FileUtils
}

func NewPDFProcessor(fileUtils *utils.FileUtils) *PDFProcessor {
	return &PDFProcessor{
		FileUtils: fileUtils,
	}
}

func (p *PDFProcessor) pdfExtension(file string) string {
	return file + ".pdf"
}

func (p *PDFProcessor) validateInputFiles(inputFilesNames []string) error {
	if len(inputFilesNames) == 0 {
		return errors.New("no PDF files provided")
	} else if len(inputFilesNames) == 1 {
		return errors.New("please provide more than one file to merge pdfs")
	}
	return nil
}

func (p *PDFProcessor) MergePdfs(pdfs []string, outputPdf string) error {
	outputPdf = p.pdfExtension(outputPdf)
	if err := p.validateInputFiles(pdfs); err != nil {
		return err
	}

	if err := api.MergeCreateFile(pdfs, outputPdf, false, nil); err != nil {
		return err
	}

	if err := api.ValidateFile(outputPdf, nil); err != nil {
		return err
	}
	return nil
}
