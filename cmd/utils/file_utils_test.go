package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createValidPDF(filepath string) error {
	// Create a minimal valid PDF file
	content := `%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Kids [3 0 R] /Count 1 >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>
endobj
xref
0 4
0000000000 65535 f 
0000000010 00000 n 
0000000053 00000 n 
0000000102 00000 n 
trailer
<< /Root 1 0 R /Size 4 >>
startxref
150
%%EOF`
	return os.WriteFile(filepath, []byte(content), 0644)
}

func createTestFiles(t *testing.T, tempDir string, pdfs []string) {
	// Create test files
	for _, f := range pdfs {
		fullPath := filepath.Join(tempDir, f)
		if err := createValidPDF(fullPath); err != nil {
			assert.NoError(t, err, "failed to create test file: ", f)
		}
	}
	// Create a subdirectory
	subDir := "subdir"
	if err := os.Mkdir(filepath.Join(tempDir, subDir), 0755); err != nil {
		assert.NoError(t, err, "failed to create subdir: ", subDir)
	}
}

func TestIsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// Create Test Subdirectory
	testSubDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(testSubDir, 0755)
	assert.NoError(t, err)

	f := &FileUtils{}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "should return true for existing directory",
			path:     tempDir,
			expected: true,
		},
		{
			name:     "should return true for existing subdirectory",
			path:     testSubDir,
			expected: true,
		},
		{
			name:     "should return false for existing file",
			path:     testFile,
			expected: false,
		},
		{
			name:     "should return false for non-existing path",
			path:     filepath.Join(tempDir, "non-existing"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.IsDirectory(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCurrentWorkingDir(t *testing.T) {
	expectedDir, err := os.Getwd()
	assert.NoError(t, err, "failed to get current working dir: %v", err)

	fileUtils := NewFileUtils(nil)
	actualDir, err := fileUtils.GetCurrentWorkingDir()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, actualDir)
}

func TestReadDirectory(t *testing.T) {
	tempDir := t.TempDir()
	err := os.Chdir(tempDir)
	assert.NoError(t, err, "failed to change directory: ", tempDir)
	createTestFiles(t, tempDir, []string{"file1.pdf", "file2.pdf"})

	f := NewFileUtils(nil)
	entries, err := f.ReadDirectory(tempDir)
	assert.NoError(t, err, "error GetCurrentWorkingDir: %v", err)

	assert.Len(t, entries, 3)

	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	assert.Contains(t, fileNames, "file1.pdf")
	assert.Contains(t, fileNames, "file2.pdf")
}

func TestFilterPdfFiles(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(t, tempDir, []string{"file1.pdf", "file2.pdf", "file3.txt"})

	f := NewFileUtils(nil)
	entries, _ := f.ReadDirectory(tempDir)

	pdfFiles := f.FilterPdfFiles(entries)
	assert.Len(t, pdfFiles, 2)
	assert.Contains(t, pdfFiles, "file1.pdf")
	assert.Contains(t, pdfFiles, "file2.pdf")
}

func TestGetPdfFilesFromDir(t *testing.T) {
	tests := []struct {
		name          string
		pdfs          []string
		dir           string
		NumberOfItems int
		expectedErr   bool
	}{
		{
			name:          "Get 2 files",
			pdfs:          []string{"file1.pdf", "file2.pdf"},
			dir:           "tempDir",
			NumberOfItems: 2,
			expectedErr:   false,
		},
		{
			name:          "fake directory",
			pdfs:          nil,
			dir:           "falseDir",
			NumberOfItems: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)
			createTestFiles(t, tempDir, tt.pdfs)

			f := NewFileUtils(nil)
			var pdfFiles []string
			if tt.dir == "tempDir" {
				pdfFiles, err = f.GetPdfFilesFromDir(tempDir)
			} else {
				pdfFiles, err = f.GetPdfFilesFromDir(tt.dir)
			}

			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but command ran successfully")
			} else {
				assert.NoError(t, err, "Expected to run successfully but received an error")
			}

			assert.Len(t, pdfFiles, tt.NumberOfItems)
			assert.Equal(t, pdfFiles, tt.pdfs)
		})
	}
}

func TestAddFullPathToPdfs(t *testing.T) {
	fileUtils := NewFileUtils(nil)

	dir := "testDir"
	pdfs := []string{"file1.pdf", "file2.pdf"}
	fullPaths := fileUtils.AddFullPathToPdfs(dir, pdfs)

	assert.Len(t, fullPaths, 2)
	assert.Equal(t, filepath.Join(dir, "file1.pdf"), fullPaths[0])
	assert.Equal(t, filepath.Join(dir, "file2.pdf"), fullPaths[1])
}

func TestCheckProvidedArgs(t *testing.T) {
	tests := []struct {
		name                string
		setup               func(tempDir string)
		args                []string
		expectedPdfs        []string
		expectedErr         bool
		expectedInteractive bool
	}{
		{
			name: "no args provided",
			setup: func(tempDir string) {
				file1 := filepath.Join(tempDir, "file1.pdf")
				file2 := filepath.Join(tempDir, "file2.pdf")
				err := createValidPDF(file1)
				assert.NoError(t, err)
				err = createValidPDF(file2)
				assert.NoError(t, err)
			},
			args:                []string{},
			expectedPdfs:        []string{"file1.pdf", "file2.pdf"},
			expectedErr:         false,
			expectedInteractive: true,
		},
		{
			name: "1 directory provided",
			setup: func(tempDir string) {
				err := os.Mkdir(filepath.Join(tempDir, "test"), 0755)
				assert.NoError(t, err)
				file1 := filepath.Join(filepath.Join(tempDir, "test"), "file1.pdf")
				err = createValidPDF(file1)
				assert.NoError(t, err)
			},
			args:                []string{"test"},
			expectedPdfs:        []string{"file1.pdf"},
			expectedErr:         false,
			expectedInteractive: true,
		},
		{
			name: "1 file provided",
			setup: func(tempDir string) {
				file1 := filepath.Join(tempDir, "file1.pdf")
				err := createValidPDF(file1)
				assert.NoError(t, err)
			},
			args:                []string{"file1.pdf"},
			expectedPdfs:        []string{"file1.pdf"},
			expectedErr:         false,
			expectedInteractive: false,
		},
		{
			name: "invalid file provided",
			setup: func(tempDir string) {
			},
			args:                []string{"file1.pdf"},
			expectedPdfs:        nil,
			expectedErr:         true,
			expectedInteractive: false,
		},
		{
			name: "file and directory provided",
			setup: func(tempDir string) {
				file1 := filepath.Join(tempDir, "file1.pdf")
				err := createValidPDF(file1)
				assert.NoError(t, err)
				err = os.Mkdir(filepath.Join(tempDir, "test"), 0755)
				assert.NoError(t, err)
			},
			args:                []string{"file1.pdf", "test"},
			expectedPdfs:        nil,
			expectedErr:         true,
			expectedInteractive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err)

			if tt.setup != nil {
				tt.setup(tempDir)
			}

			f := NewFileUtils(nil)
			f.args = tt.args
			pdfs, _, err := f.CheckProvidedArgs()

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPdfs, pdfs)

			switch tt.expectedInteractive {
			case true:
				assert.True(t, f.Interactive, "expected interactive to be true")
			case false:
				assert.False(t, f.Interactive, "expected interactive to be false")
			}
		})
	}
}
