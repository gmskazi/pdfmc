package utils

import (
	"os"
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

func TestMergedPdfs(t *testing.T) {
	// Create temporary PDF files
	// By using t.TempDir() ensures all temporary files are cleaned up.
	temDir := t.TempDir()
	inputFile1 := temDir + "/file1.pdf"
	inputFile2 := temDir + "/file2.pdf"
	outputfile := temDir + "/merged.pdf"

	// Create dummy pdfs
	err := createValidPDF(inputFile1)
	assert.NoError(t, err)
	err = createValidPDF(inputFile2)
	assert.NoError(t, err)

	err = MergePdfs([]string{inputFile1, inputFile2}, outputfile)
	assert.NoError(t, err)

	_, err = os.Stat(outputfile)
	assert.NoError(t, err)
}

func TestValidateInputFiles(t *testing.T) {
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

			pdfs, _ := GetPdfFilesFromDir(tempDir)

			err := validateInputFiles(pdfs)

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
	actual := PdfExtension("testing")

	if actual != expected {
		t.Errorf("Expected: %s got %s", expected, actual)
	}
}

func TestGetCurrentWorkingDir(t *testing.T) {
	expectedDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working dir: %v", err)
	}

	actualDir, err := GetCurrentWorkingDir()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, actualDir)
}
