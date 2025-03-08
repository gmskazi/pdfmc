package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func GetCurrentWorkingDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working dir: %w", err)
	}
	return dir, nil
}

func ReadDirectory(directory string) ([]os.DirEntry, error) {
	return os.ReadDir(directory)
}

func FilterPdfFiles(directory string, entries []os.DirEntry) ([]string, error) {
	var pdfFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, entry.Name())
		}
	}

	return pdfFiles, nil
}

func GetPdfFilesFromDir(directory string) ([]string, error) {
	entries, err := ReadDirectory(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s", err)
	}

	return FilterPdfFiles(directory, entries)
}

func PdfExtension(file string) string {
	return file + ".pdf"
}

func validateInputFiles(inputFilesNames []string) error {
	if len(inputFilesNames) == 0 {
		return errors.New("no PDF files provided")
	} else if len(inputFilesNames) == 1 {
		return errors.New("please provide more than one file to merge pdfs")
	}
	return nil
}

func MergePdfs(pdfs []string, outputPdf string) error {
	if err := validateInputFiles(pdfs); err != nil {
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
