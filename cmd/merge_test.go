package cmd

import (
	"bytes"
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

// Only testing non interactive mode for now
func TestMergeCommand(t *testing.T) {
	file1 := "file1.pdf"
	file2 := "file2.pdf"

	tests := []struct {
		name           string
		pdfs           []string
		flags          []string
		fileOutput     string
		expectError    bool
		expectedOutput string
		checkFile      bool
	}{
		{
			name:           "Merge two PDF files",
			pdfs:           []string{file1, file2},
			flags:          []string{merge, file1, file2},
			fileOutput:     "merged_output.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
		{
			name:           "Merge two PDF files with custom filename",
			pdfs:           []string{file1, file2},
			flags:          []string{merge, file1, file2, "-n", "testname"},
			fileOutput:     "testname.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
		{
			name:           "Merge two PDF files with custom filename and password",
			pdfs:           []string{file1, file2},
			flags:          []string{merge, file1, file2, "-n", "testname", "-p", "test"},
			fileOutput:     "encrypted-testname.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged and encrypted successfully to:",
			checkFile:      true,
		},

		{
			name:           "Check if file and directory is provided",
			pdfs:           []string{file1},
			flags:          []string{merge, file1, "subdir"},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "is a directory not a pdf",
			checkFile:      false,
		},
		{
			name:           "Check if provided files are avalible.",
			pdfs:           nil,
			flags:          []string{merge, file1, file2},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "no such file or directory",
			checkFile:      false,
		},
		{
			name:           "Check if directory is valid.",
			pdfs:           nil,
			flags:          []string{merge, "tempDir"},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "no such file or directory",
			checkFile:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)

			createTestFiles(t, tempDir, tt.pdfs)
			args := tt.flags

			var outputBuf bytes.Buffer

			rootCmd.SetOut(&outputBuf)
			rootCmd.SetErr(&outputBuf)
			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if tt.expectError {
				assert.Error(t, err, "Expected an error but command ran successfuly.")
			} else {
				assert.NoError(t, err, "Expected command to run successfuly but it failed.")
			}

			assert.Contains(t, outputBuf.String(), tt.expectedOutput, "Unexpected output from command.")
			// fmt.Println(outputBuf.String())

			if tt.checkFile {
				_, err := os.Stat(tt.fileOutput)
				assert.NoError(t, err, "Expected merged PDF file to be created but it wasn't there.")
			}
		})
	}
}
