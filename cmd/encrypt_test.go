package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Only testing non interactive mode for now
func TestEncryptCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		setup          func(t *testing.T, tempDir string) []string
		fileOutput     string
		expectError    bool
		expectedOutput string
		checkFile      bool
	}{
		{
			name: "Encrypt a PDF file",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				return []string{"encrypt", file1, "-p", "test"}
			},
			fileOutput:     "encrypted-file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file encrypted successfully to:",
			checkFile:      true,
		},
		{
			name: "Encrypt multiple PDF files",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"
				file2 := "file2.pdf"
				file3 := "file3.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(filepath.Join(tempDir, file2))
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}
				err = createValidPDF(filepath.Join(tempDir, file3))
				if err != nil {
					t.Fatalf("failed to create file3.pdf: %v", err)
				}
				return []string{"encrypt", file1, file2, file3, "-p", "test"}
			},
			fileOutput:     "encrypted-file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file encrypted successfully to:",
			checkFile:      true,
		},
		{
			name: "Provide a flag with no password",
			setup: func(t *testing.T, tempDir string) []string {
				file1 := "file1.pdf"

				err := createValidPDF(filepath.Join(tempDir, file1))
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				return []string{"encrypt", file1, "-p"}
			},
			fileOutput:     "",
			expectError:    true,
			expectedOutput: "flag needs an argument:",
			checkFile:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)

			args := tt.setup(t, tempDir)

			var outputBuf bytes.Buffer

			rootCmd.SetOut(&outputBuf)
			rootCmd.SetErr(&outputBuf)
			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if tt.expectError {
				assert.Error(t, err, "Expected an error but command ran successfully.")
			} else {
				assert.NoError(t, err, "Expected command to run successfully but it failed.")
			}

			assert.Contains(t, outputBuf.String(), tt.expectedOutput, "Expected output to contain: %s", tt.expectedOutput)

			if tt.checkFile {
				_, err := os.Stat(tt.fileOutput)
				assert.NoError(t, err, "Expected file %s to be created but it was not found.", tt.fileOutput)
			}
		})
	}
}
