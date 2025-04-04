package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gmskazi/pdfmc/cmd/utils"
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

func TestMergedPdfs(t *testing.T) {
	// Create temporary PDF files
	// By using t.TempDir() ensures all temporary files are cleaned up.
	temDir := t.TempDir()
	inputFile1 := filepath.Join(temDir, "file1.pdf")
	inputFile2 := filepath.Join(temDir, "/file2.pdf")
	customName := "merged"
	outputfile := "merged.pdf"

	// Create dummy pdfs
	err := createValidPDF(inputFile1)
	assert.NoError(t, err)
	err = createValidPDF(inputFile2)
	assert.NoError(t, err)

	fileUtils := utils.NewFileUtils(nil)
	pdfProcssor := NewPDFProcessor(fileUtils, "merge")
	output, err := pdfProcssor.MergePdfs([]string{inputFile1, inputFile2}, customName)
	assert.NoError(t, err)

	assert.Equal(t, output, outputfile)
	_, err = os.Stat(outputfile)
	assert.NoError(t, err)
}

func TestValidateInputFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setup         func(tempDir string)
		expectedFiles []string
		expectedErr   string
	}{
		{
			name: "No PDF files are provided",
			setup: func(tempDir string) {
			},
			expectedFiles: nil,
			expectedErr:   "no PDF files provided",
		},
		{
			name: "One PDF file provided",
			setup: func(tempDir string) {
				err := createValidPDF(tempDir + "/file1.pdf")
				if err != nil {
					t.Fatalf("failed to write file1.pdf: %v", err)
				}
			},
			expectedFiles: []string{"file1.pdf"},
			expectedErr:   "please provide more than one file to merge pdfs",
		},
		{
			name: "Multiple PDF files to merge",
			setup: func(temDir string) {
				err := createValidPDF(temDir + "/file1.pdf")
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(temDir + "/file2.pdf")
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}
				err = createValidPDF(temDir + "/file3.pdf")
				if err != nil {
					t.Fatalf("failed to create file3.pdf: %v", err)
				}
			},
			expectedFiles: []string{"file1.pdf", "file2.pdf", "file3.pdf"},
			expectedErr:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.setup != nil {
				tt.setup(tempDir)
			}

			fileUtils := utils.NewFileUtils(nil)
			pdfProcessor := NewPDFProcessor(fileUtils, "merge")

			pdfs := fileUtils.GetPdfFilesFromDir(tempDir)
			fmt.Println(pdfs)

			err := pdfProcessor.validateInputFiles(pdfs)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, pdfs, tt.expectedFiles)
			}
		})
	}
}

func TestPdfExtension(t *testing.T) {
	expected := "testing.pdf"
	fileUtils := utils.NewFileUtils(nil)
	pdfProcessor := NewPDFProcessor(fileUtils, "merge")
	actual := pdfProcessor.pdfExtension("testing")

	if actual != expected {
		t.Errorf("Expected: %s got %s", expected, actual)
	}
}

func TestEncryptPdf(t *testing.T) {
	tests := []struct {
		name         string
		pdf          string
		password     string
		expectedFile string
		expectedErr  bool
		setupFile    bool
	}{
		{
			name:         "successful encryption",
			pdf:          "test.pdf",
			password:     "test",
			expectedFile: "encrypted-test.pdf",
			expectedErr:  false,
			setupFile:    true,
		},
	}

	for _, tt := range tests {
		tempDir := t.TempDir()

		if tt.setupFile {
			testPdf := filepath.Join(tempDir, tt.pdf)
			err := createValidPDF(testPdf)
			if err != nil {
				t.Fatalf("failed to create test.pdf: %v", err)
			}

			processor := NewPDFProcessor(utils.NewFileUtils(nil), "encrypt")
			encryptedPdf, err := processor.EncryptPdf(tt.pdf, tempDir, tt.password)
			if err != nil {
				t.Fatalf("failed to encrypt PDF: %v", err)
			}

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Empty(t, encryptedPdf)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedFile, encryptedPdf)

		}
	}
}
