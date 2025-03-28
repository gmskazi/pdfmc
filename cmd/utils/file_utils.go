package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileUtils struct {
	pdfs        []string
	interactive bool
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
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get currect working dir: %w", err)
	}

	return dir, nil
}

func (f *FileUtils) ReadDirectory(directory string) ([]os.DirEntry, error) {
	return os.ReadDir(directory)
}

func (f *FileUtils) FilterPdfFiles(directory string, entries []os.DirEntry) []string {
	var pdfFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, entry.Name())
		}
	}

	return pdfFiles
}

func (f *FileUtils) GetPdfFilesFromDir(directory string) []string {
	entries, _ := f.ReadDirectory(directory)

	return f.FilterPdfFiles(directory, entries)
}

func (f *FileUtils) AddFullPathToPdfs(dir string, pdfs []string) []string {
	var fullPaths []string

	for _, pdf := range pdfs {
		fullPaths = append(fullPaths, filepath.Join(dir, pdf))
	}
	return fullPaths
}

func (f *FileUtils) CheckProvidedArgs(args []string) (pdfs []string, interactive bool, err error) {
	// check if any files/folders are provided
	if len(args) == 0 {
		f.interactive = true
		f.dir, err = f.GetCurrentWorkingDir()
		if err != nil {
			return []string{}, true, err
		}
		f.pdfs = f.GetPdfFilesFromDir(f.dir)
		return f.pdfs, f.interactive, nil

	} else if len(args) == 1 && f.IsDirectory(args[0]) {
		f.interactive = true
		f.dir = args[0]
		f.pdfs = f.GetPdfFilesFromDir(f.dir)
		return f.pdfs, f.interactive, nil

	} else {
		f.interactive = false
		f.pdfs = args
		for _, pdf := range f.pdfs {
			if _, err := os.Stat(pdf); os.IsNotExist(err) {
				return []string{}, false, err
			}
			if info, err := os.Stat(pdf); err == nil && info.IsDir() {
				return []string{}, false, err
			}
		}
	}
	return f.pdfs, f.interactive, nil
}
