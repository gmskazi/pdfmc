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

// Only testing non interactive mode for now
func TestMergeCommand(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T, tempDir string) []string
		fileOutput     string
		expectError    bool
		expectedOutput string
		checkFile      bool
	}{
		{
			name: "Merge two PDF files",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				file2 := "file2.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(filepath.Join(tempDir, file2))
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}

				return []string{"merge", file1, file2}
			},
			fileOutput:     "merged_output.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
		{
			name: "Merge two PDF files with custom filename",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				file2 := "file2.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(filepath.Join(tempDir, file2))
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}

				return []string{"merge", file1, file2, "-n", "testname"}
			},
			fileOutput:     "testname.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
		{
			name: "Merge two PDF files with custom filename and password",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				file2 := "file2.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				assert.NoError(t, err, "failed to create file1.pdf: %v", err)
				err = createValidPDF(filepath.Join(tempDir, file2))
				assert.NoError(t, err, "failed to create file2.pdf: %v", err)

				return []string{"merge", file1, file2, "-n", "testname", "-p", "test"}
			},
			fileOutput:     "encrypted-testname.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},

		{
			name: "Check if file and directory is provided",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				dir := filepath.Join(tempDir, "test")

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = os.Mkdir(dir, 0755)
				if err != nil {
					t.Fatalf("failed to create test directory: %v", err)
				}

				return []string{"merge", file1, dir}
			},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "is a directory not a pdf",
			checkFile:      false,
		},
		{
			name: "Check if provided files are avalible.",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				file2 := "file2.pdf"

				return []string{"merge", file1, file2}
			},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "no such file or directory",
			checkFile:      false,
		},
		{
			name: "Check if directory is valid.",
			setup: func(t *testing.T, tempDir string) []string {
				err := os.RemoveAll(tempDir)
				assert.NoError(t, err, "failed to remove directory: ", tempDir)

				return []string{"merge", tempDir}
			},
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

			args := tt.setup(t, tempDir)
			// fmt.Println(args)

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
