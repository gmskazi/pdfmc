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
	if err != nil {
		t.Fatalf("failed to get current working dir: %v", err)
	}

	fileUtils := *NewFileUtils(nil)
	actualDir, err := fileUtils.GetCurrentWorkingDir()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, actualDir)

	fileUtils.dir = "dummydir"
	_, err = fileUtils.GetCurrentWorkingDir()
	assert.NoError(t, err)
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

	fileUtils := NewFileUtils(nil)
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
	fileUtils := NewFileUtils(nil)

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

	pdfFiles := fileUtils.FilterPdfFiles(entries)
	assert.Len(t, pdfFiles, 2)
	assert.Contains(t, pdfFiles, "file1.pdf")
	assert.Contains(t, pdfFiles, "file2.pdf")
}

func TestGetPdfFilesFromDir(t *testing.T) {
	fileUtils := NewFileUtils(nil)

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
	fileUtils := NewFileUtils(nil)

	dir := "testDir"
	pdfs := []string{"file1.pdf", "file2.pdf"}
	fullPaths := fileUtils.AddFullPathToPdfs(dir, pdfs)

	assert.Len(t, fullPaths, 2)
	assert.Equal(t, filepath.Join(dir, "file1.pdf"), fullPaths[0])
	assert.Equal(t, filepath.Join(dir, "file2.pdf"), fullPaths[1])
}

func TestCheckProvidedArgs(t *testing.T) {
	t.Parallel()

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

			fileUtils := NewFileUtils(nil)
			fileUtils.args = tt.args
			pdfs, interactive, err := fileUtils.CheckProvidedArgs()

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPdfs, pdfs)

			switch tt.expectedInteractive {
			case true:
				assert.True(t, interactive, "expected interactive to be true")
			case false:
				assert.False(t, interactive, "expected interactive to be false")
			}
		})
	}
}
