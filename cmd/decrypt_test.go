package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/stretchr/testify/assert"
)

func encryptTestFiles(t *testing.T, tempdir string, pdfs []string, password, pdfPrefix string) {
	p := pdf.NewPDFProcessor(encrypt)
	// encrypt test files
	for _, f := range pdfs {

		encryptedPdf, err := p.EncryptPdf(f, tempdir, password, pdfPrefix)
		assert.NoError(t, err, "failed to encrypt pdf")
		fmt.Println(encryptedPdf)

	}
}

// Only testing non interactive mode for now
func TestDecryptCommand(t *testing.T) {
	tests := []struct {
		name           string
		pdfs           []string
		pdfPrefix      string
		flags          []string
		fileOutput     string
		expectError    bool
		expectedOutput string
		checkFile      bool
		encrypt        bool
		password       string
	}{
		{
			name:           "Decrypt a PDF file",
			pdfs:           []string{"file1.pdf"},
			pdfPrefix:      "",
			flags:          []string{decrypt, "file1.pdf", "-p", "test"},
			fileOutput:     "file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file decrypted successfully to:",
			checkFile:      true,
			encrypt:        true,
			password:       "test",
		},
		{
			name:           "Check if file is not encrypted",
			pdfs:           []string{"file1.pdf"},
			flags:          []string{decrypt, "file1.pdf", "-p", "test"},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "this file is not encrypted",
			checkFile:      false,
			encrypt:        false,
		},
		{
			name:           "Provide a flag with no password",
			pdfs:           nil,
			pdfPrefix:      "",
			flags:          []string{decrypt, "file1.pdf", "-p"},
			fileOutput:     "",
			expectError:    true,
			expectedOutput: "flag needs an argument:",
			checkFile:      false,
			encrypt:        false,
			password:       "",
		},
		{
			name:           "Check if files are available",
			pdfs:           nil,
			pdfPrefix:      "",
			flags:          []string{decrypt, "file1.pdf", "-p", "test"},
			fileOutput:     "",
			expectError:    false,
			expectedOutput: "no such file or directory",
			checkFile:      false,
			encrypt:        false,
			password:       "",
		},
		{
			name:           "decrypt file with custom name prefix",
			pdfs:           []string{"file1.pdf"},
			pdfPrefix:      "",
			flags:          []string{decrypt, "file1.pdf", "-p", "test", "-n", "test-"},
			fileOutput:     "test-file1.pdf",
			expectError:    false,
			expectedOutput: "PDF file decrypted successfully to:",
			checkFile:      true,
			encrypt:        true,
			password:       "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err := os.Chdir(tempDir)
			assert.NoError(t, err, "failed to change directory: ", tempDir)

			createTestFiles(t, tempDir, tt.pdfs)
			if tt.encrypt && tt.password != "" {
				encryptTestFiles(t, tempDir, tt.pdfs, tt.password, tt.pdfPrefix)
			}
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
