package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Only testing non interactive mode for now
func TestEncryptCommand(t *testing.T) {
	t.Parallel()
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
			name:           "Encrypt a PDF file",
			pdfs:           []string{"file1.pdf"},
			flags:          []string{encrypt, "file1.pdf", "-p", "test"},
			fileOutput:     "encrypted-file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file encrypted successfully to:",
			checkFile:      true,
		},
		{
			name:           "Encrypt multiple PDF files",
			pdfs:           []string{"file1.pdf", "file2.pdf", "file3.pdf"},
			flags:          []string{encrypt, "file1.pdf", "file2.pdf", "file3.pdf", "-p", "test"},
			fileOutput:     "encrypted-file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file encrypted successfully to:",
			checkFile:      true,
		},
		{
			name:           "Provide a flag with no password",
			pdfs:           nil,
			flags:          []string{encrypt, "file1.pdf", "-p"},
			fileOutput:     "",
			expectError:    true,
			expectedOutput: "flag needs an argument:",
			checkFile:      false,
		},
		{
			name:           "Check if files are available",
			pdfs:           nil,
			flags:          []string{encrypt, "file1.pdf", "-p", "test"},
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
