package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileUtils struct {
	pdfs        []string
	Interactive bool
	dir         string
	args        []string
}

func NewFileUtils(args []string) *FileUtils {
	return &FileUtils{
		args: args,
	}
}

func (f *FileUtils) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (f *FileUtils) GetCurrentWorkingDir() (string, error) {
	var err error
	f.dir, err = os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get currect working dir: %w", err)
	}

	return f.dir, nil
}

func (f *FileUtils) ReadDirectory(directory string) ([]os.DirEntry, error) {
	return os.ReadDir(directory)
}

func (f *FileUtils) FilterPdfFiles(entries []os.DirEntry) []string {
	var pdfFiles []string

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, entry.Name())
		}
	}
	f.pdfs = pdfFiles

	return f.pdfs
}

func (f *FileUtils) GetPdfFilesFromDir(directory string) ([]string, error) {
	entries, err := f.ReadDirectory(directory)
	if err != nil {
		return nil, err
	}

	pdfFiles := f.FilterPdfFiles(entries)
	return pdfFiles, nil
}

func (f *FileUtils) AddFullPathToPdfs(dir string, pdfs []string) []string {
	var fullPaths []string

	for _, pdf := range pdfs {
		fullPaths = append(fullPaths, filepath.Join(dir, pdf))
	}
	return fullPaths
}

func (f *FileUtils) CheckProvidedArgs() ([]string, string, error) {
	var err error

	if len(f.args) == 0 {
		f.Interactive = true
		f.dir, err = f.GetCurrentWorkingDir()
		if err != nil {
			return nil, f.dir, err
		}
		f.pdfs, err = f.GetPdfFilesFromDir(f.dir)
		if err != nil {
			return nil, f.dir, err
		}
		return f.pdfs, f.dir, nil

	}

	if len(f.args) == 1 && f.IsDirectory(f.args[0]) {
		f.Interactive = true
		f.dir = f.args[0]
		f.pdfs, err = f.GetPdfFilesFromDir(f.dir)
		if err != nil {
			return nil, f.dir, err
		}
		return f.pdfs, f.dir, nil

	}

	f.Interactive = false
	for _, pdf := range f.args {
		info, err := os.Stat(pdf)
		if err != nil {
			return nil, "", err
		}
		if info.IsDir() {
			return nil, "", fmt.Errorf("%s is a directory not a pdf", pdf)
		}
	}

	f.pdfs = f.args
	return f.pdfs, f.dir, nil
}
