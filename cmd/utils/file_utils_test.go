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

func TestGetCurrentWorkingDir(t *testing.T) {
	expectedDir, _ := os.Getwd()

	fileUtils := NewFileUtils()
	actualDir, err := fileUtils.GetCurrentWorkingDir()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, actualDir)
}

func TestReadDirectory(t *testing.T) {
	tempDir := t.TempDir()
	inputFile1 := tempDir + "/file1.pdf"
	inputFile2 := tempDir + "/file2.pdf"
	err := createValidPDF(inputFile1)
	if err != nil {
		t.Errorf("error creating file1.pdf: %v", err)
	}
	err = createValidPDF(inputFile2)
	if err != nil {
		t.Errorf("error createing file2.pdf: %v", err)
	}

	fileUtils := NewFileUtils()
	entries, err := fileUtils.ReadDirectory(tempDir)
	if err != nil {
		t.Errorf("error GetCurrentWorkingDir: %v", err)
	}

	assert.NoError(t, err)
	assert.Len(t, entries, 2)

	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	assert.Contains(t, fileNames, "file1.pdf")
	assert.Contains(t, fileNames, "file2.pdf")
}

func TestFilterPdfFiles(t *testing.T) {
	fileUtils := NewFileUtils()

	tempDir := t.TempDir()
	err := createValidPDF(tempDir + "/file1.pdf")
	if err != nil {
		t.Errorf("error creating file1.pdf: %v", err)
	}
	err = createValidPDF(tempDir + "/file2.pdf")
	if err != nil {
		t.Errorf("error creating file2.pdf: %v", err)
	}
	err = createValidPDF(tempDir + "/file1.txt")
	if err != nil {
		t.Errorf("error creating file1.txt: %v", err)
	}

	err = os.Mkdir(tempDir+"/testing", 0755)
	if err != nil {
		t.Errorf("error creating testing dir: %v", err)
	}

	entries, _ := fileUtils.ReadDirectory(tempDir)

	pdfFiles := fileUtils.FilterPdfFiles("", entries)
	assert.Len(t, pdfFiles, 2)
	assert.Contains(t, pdfFiles, "file1.pdf")
	assert.Contains(t, pdfFiles, "file2.pdf")
}

func TestGetPdfFilesFromDir(t *testing.T) {
	fileUtils := NewFileUtils()

	tempDir := t.TempDir()
	err := createValidPDF(tempDir + "/file1.pdf")
	if err != nil {
		t.Errorf("error creating file1.pdf: %v", err)
	}
	err = createValidPDF(tempDir + "/file2.pdf")
	if err != nil {
		t.Errorf("error creating file2.pdf: %v", err)
	}
	pdfFiles := fileUtils.GetPdfFilesFromDir(tempDir)

	assert.Len(t, pdfFiles, 2)
	assert.Contains(t, pdfFiles, "file1.pdf")
	assert.Contains(t, pdfFiles, "file2.pdf")
}

func TestAddFullPathToPdfs(t *testing.T) {
	fileUtils := NewFileUtils()

	dir := "testDir"
	pdfs := []string{"file1.pdf", "file2.pdf"}
	fullPaths := fileUtils.AddFullPathToPdfs(dir, pdfs)

	assert.Len(t, fullPaths, 2)
	assert.Equal(t, filepath.Join(dir, "file1.pdf"), fullPaths[0])
	assert.Equal(t, filepath.Join(dir, "file2.pdf"), fullPaths[1])
}
