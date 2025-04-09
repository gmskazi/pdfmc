package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	merge   = "merge"
	encrypt = "encrypt"
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

func TestMergePdfs(t *testing.T) {
	file1 := "file1.pdf"
	file2 := "file2.pdf"

	tests := []struct {
		name        string
		pdfs        []string
		customName  string
		fileOutput  string
		expectedErr bool
	}{
		{
			name:        "Merge 2 files",
			pdfs:        []string{file1, file2},
			customName:  "test",
			fileOutput:  "test.pdf",
			expectedErr: false,
		},
		{
			name:        "Error for one file",
			pdfs:        []string{file1},
			customName:  "test",
			fileOutput:  "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)

			createTestFiles(t, tempDir, tt.pdfs)
			pdfP := NewPDFProcessor(merge)
			output, err := pdfP.MergePdfs(tt.pdfs, tt.customName)

			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but command ran successfully")
			} else {
				assert.NoError(t, err, "Expected command to run successfully but it failed")
			}

			assert.Equal(t, tt.fileOutput, output, "File output %s should be the same as %s", output, tt.fileOutput)
		})
	}
}

func TestPdfExtension(t *testing.T) {
	expected := "testing.pdf"
	pdfProcessor := NewPDFProcessor("merge")
	actual := pdfProcessor.pdfExtension("testing")

	assert.Equal(t, expected, actual, "Expected: %s, got %s", expected, actual)
}

func TestEncryptPdf(t *testing.T) {
	tests := []struct {
		name         string
		pdf          string
		password     string
		expectedFile string
		expectedErr  bool
		setupFile    []string
	}{
		{
			name:         "successful encryption",
			pdf:          "test.pdf",
			password:     "test",
			expectedFile: "encrypted-test.pdf",
			expectedErr:  false,
			setupFile:    []string{"test.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: %s", tempDir)
			createTestFiles(t, tempDir, tt.setupFile)

			processor := NewPDFProcessor(encrypt)
			encryptedPdf, err := processor.EncryptPdf(tt.pdf, tempDir, tt.password)
			fmt.Println(err)
			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but command ran successfully")
			} else {
				assert.NoError(t, err, "Expected to run successfully but it failed")
			}

			assert.Equal(t, tt.expectedFile, encryptedPdf, "Expected PDF file: %s to be Equal to: %s", encryptedPdf, tt.expectedFile)
		})
	}
}
