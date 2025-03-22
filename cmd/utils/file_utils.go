package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileUtils struct{}

func NewFileUtils() *FileUtils {
	return &FileUtils{}
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
