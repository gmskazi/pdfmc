package pdf

import (
	"errors"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFProcessor struct {
	logo string
}

func NewPDFProcessor(logo string) *PDFProcessor {
	return &PDFProcessor{
		logo: logo,
	}
}

func (p *PDFProcessor) pdfExtension(file string) string {
	if filepath.Ext(file) != ".pdf" {
		return file + ".pdf"
	}
	return file
}

func (p *PDFProcessor) MergePdfs(pdfs []string, outputPdf string) (string, error) {
	if len(pdfs) < 2 {
		return "", errors.New("at least two PDF files are required to merge")
	}
	output := p.pdfExtension(outputPdf)

	if err := api.MergeCreateFile(pdfs, output, false, nil); err != nil {
		return "", err
	}

	if err := api.ValidateFile(output, nil); err != nil {
		return "", err
	}
	return output, nil
}

func (p *PDFProcessor) EncryptPdf(pdf, dir, password, prefix string) (string, error) {
	var encryptedPdfName string
	conf := model.NewAESConfiguration(password, password, 256)
	if prefix != "" {
		encryptedPdfName = prefix + pdf
	} else {
		encryptedPdfName = pdf
	}

	err := api.EncryptFile(filepath.Join(dir, pdf), "", conf)
	if err != nil {
		return "", err
	}

	err = api.ValidateFile(pdf, conf)
	if err != nil {
		return encryptedPdfName, err
	}
	return encryptedPdfName, nil
}

func (p *PDFProcessor) DecryptPdf(pdf, dir, password, prefix string) (string, error) {
	var decryptedPdfName string
	conf := model.NewAESConfiguration(password, password, 256)
	if prefix != "" {
		decryptedPdfName = prefix + pdf
	} else {
		decryptedPdfName = pdf
	}

	err := api.DecryptFile(filepath.Join(dir, pdf), "", conf)
	if err != nil {
		return "", err
	}

	err = api.ValidateFile(decryptedPdfName, conf)
	if err != nil {
		return "", err
	}
	return decryptedPdfName, nil
}
