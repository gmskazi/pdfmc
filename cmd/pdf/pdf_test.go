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
	decrypt = "decrypt"
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
		err := createValidPDF(fullPath)
		assert.NoError(t, err, "failed to create test file: ", f)
	}
	// Create a subdirectory
	subDir := "subdir"
	err := os.Mkdir(filepath.Join(tempDir, subDir), 0755)
	assert.NoError(t, err, "failed to create subdir: ", subDir)
}

func encryptTestFiles(t *testing.T, tempdir, pdf, password, pdfPrefix string) {
	p := NewPDFProcessor(encrypt)
	// encrypt test files
	encryptedPdf, err := p.EncryptPdf(pdf, tempdir, password, pdfPrefix)
	assert.NoError(t, err, "failed to encrypt pdf")
	// NOTE: This is for debuging
	fmt.Println(encryptedPdf)
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
		setupFile   []string
	}{
		{
			name:        "Merge 2 files",
			pdfs:        []string{file1, file2},
			customName:  "test",
			fileOutput:  "test.pdf",
			expectedErr: false,
			setupFile:   []string{file1, file2},
		},
		{
			name:        "Error for one file",
			pdfs:        []string{file1},
			customName:  "test",
			fileOutput:  "",
			expectedErr: true,
			setupFile:   []string{file1},
		},
		{
			name:        "No files provided",
			pdfs:        []string{file1, file2},
			customName:  "test",
			fileOutput:  "",
			expectedErr: true,
			setupFile:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)

			createTestFiles(t, tempDir, tt.setupFile)
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
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name:     "add pdf extension",
			file:     "testing",
			expected: "testing.pdf",
		},
		{
			name:     "file with pdf extension",
			file:     "testing.pdf",
			expected: "testing.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfProcessor := NewPDFProcessor(merge)
			actual := pdfProcessor.pdfExtension(tt.file)

			assert.Equal(t, tt.expected, actual, "Expected: %s, got %s", tt.expected, actual)
		})
	}
}

func TestEncryptPdf(t *testing.T) {
	tests := []struct {
		name         string
		pdf          string
		pdfPrefix    string
		password     string
		expectedFile string
		expectedErr  bool
		setupFile    []string
	}{
		{
			name:         "successful encryption",
			pdf:          "test.pdf",
			pdfPrefix:    "",
			password:     "test",
			expectedFile: "test.pdf",
			expectedErr:  false,
			setupFile:    []string{"test.pdf"},
		},
		{
			name:         "successful encryption with custom prefix",
			pdf:          "test.pdf",
			pdfPrefix:    "encrypt-",
			password:     "test",
			expectedFile: "encrypt-test.pdf",
			expectedErr:  false,
			setupFile:    []string{"test.pdf"},
		},
		{
			name:         "No file provided",
			pdf:          "",
			pdfPrefix:    "",
			password:     "",
			expectedFile: "",
			expectedErr:  true,
			setupFile:    nil,
		},
		{
			name:         "No file with prefix",
			pdf:          "",
			pdfPrefix:    "test-",
			password:     "",
			expectedFile: "",
			expectedErr:  true,
			setupFile:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: %s", tempDir)
			createTestFiles(t, tempDir, tt.setupFile)

			processor := NewPDFProcessor(encrypt)
			encryptedPdf, err := processor.EncryptPdf(tt.pdf, tempDir, tt.password, tt.pdfPrefix)
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

func TestDecryptPdf(t *testing.T) {
	tests := []struct {
		name         string
		pdf          string
		pdfPrefix    string
		password     string
		expectedFile string
		expectedErr  bool
		setupFile    []string
	}{
		{
			name:         "successful decryption",
			pdf:          "test.pdf",
			pdfPrefix:    "",
			password:     "test",
			expectedFile: "test.pdf",
			expectedErr:  false,
			setupFile:    []string{"test.pdf"},
		},
		{
			name:         "successful decryption with custom prefix",
			pdf:          "test.pdf",
			pdfPrefix:    "decrypted-",
			password:     "test",
			expectedFile: "decrypted-test.pdf",
			expectedErr:  false,
			setupFile:    []string{"test.pdf"},
		},
		{
			name:         "No file provided",
			pdf:          "",
			pdfPrefix:    "",
			password:     "",
			expectedFile: "",
			expectedErr:  true,
			setupFile:    nil,
		},
		{
			name:         "No file with prefix",
			pdf:          "",
			pdfPrefix:    "test-",
			password:     "",
			expectedFile: "",
			expectedErr:  true,
			setupFile:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: %s", tempDir)
			createTestFiles(t, tempDir, tt.setupFile)
			if tt.setupFile != nil {
				encryptTestFiles(t, tempDir, tt.pdf, tt.password, "")
			}

			processor := NewPDFProcessor(decrypt)
			decryptedPdf, err := processor.DecryptPdf(tt.pdf, tempDir, tt.password, tt.pdfPrefix)
			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but command ran successfully")
			} else {
				assert.NoError(t, err, "Expected to run successfully but it failed")
			}

			assert.Equal(t, tt.expectedFile, decryptedPdf, "Expected PDF file: %s to be Equal to: %s", decryptedPdf, tt.expectedFile)
		})
	}
}
