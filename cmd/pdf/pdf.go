package pdf

import (
	"errors"
	"path/filepath"

	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFProcessor struct {
	FileUtils *utils.FileUtils
	logo      string
}

func NewPDFProcessor(fileUtils *utils.FileUtils, logo string) *PDFProcessor {
	return &PDFProcessor{
		FileUtils: fileUtils,
		logo:      logo,
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

func (p *PDFProcessor) MergePdfs(pdfs []string, outputPdf string) (output string, err error) {
	output = p.pdfExtension(outputPdf)
	if err := p.validateInputFiles(pdfs); err != nil {
		return "", err
	}

	if err := api.MergeCreateFile(pdfs, output, false, nil); err != nil {
		return "", err
	}

	if err := api.ValidateFile(output, nil); err != nil {
		return "", err
	}
	return output, nil
}

func (p *PDFProcessor) EncryptPdf(pdf string, dir string, password string) (encryptedPdf string, err error) {
	conf := model.NewAESConfiguration(password, password, 256)

	encryptedPdfName := "encrypted-" + pdf

	err = api.EncryptFile(filepath.Join(dir, pdf), encryptedPdfName, conf)
	if err != nil {
		return "", err
	}

	err = api.ValidateFile(encryptedPdfName, conf)
	if err != nil {
		return "", err
	}
	return encryptedPdfName, nil
}
