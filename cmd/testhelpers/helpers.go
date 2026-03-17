package testhelpers

import (
	"os"
	"path/filepath"
	"testing"
)

func CreateValidPDF(filepath string) error {
	content := []byte(`%PDF-1.4
1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj
2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj
3 0 obj<</Type/Page/Parent 2 0 R/MediaBox[0 0 612 792]/Resources<</Font<</F1 4 0 R>>>>/Contents 5 0 R>>endobj
4 0 obj<</Type/Font/Subtype/Type1/BaseFont/Helvetica>>endobj
5 0 obj<</Length 44>>stream
BT
/F1 12 Tf
100 700 Td
(Test PDF) Tj
ET
endstream endobj
xref
0 6
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000214 00000 n 
0000000285 00000 n 
trailer<</Size 6/Root 1 0 R>>
startxref
374
%%EOF`)
	return os.WriteFile(filepath, content, 0600)
}

func CreateTestFiles(t *testing.T, tempDir string, pdfs []string) {
	t.Helper()

	for _, pdf := range pdfs {
		err := CreateValidPDF(filepath.Join(tempDir, pdf))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", pdf, err)
		}
	}

	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	err = CreateValidPDF(filepath.Join(subDir, "nested.pdf"))
	if err != nil {
		t.Fatalf("Failed to create nested test file: %v", err)
	}
}

func EncryptTestFiles(
	t *testing.T,
	tempDir string,
	pdfs []string,
	password string,
	pdfPrefix string,
	encryptFn func(pdf, outDir, password, prefix string) (string, error),
) {
	t.Helper()

	for _, pdf := range pdfs {
		_, err := encryptFn(pdf, tempDir, password, pdfPrefix)
		if err != nil {
			t.Fatalf("Failed to encrypt test file %s: %v", pdf, err)
		}
	}
}
