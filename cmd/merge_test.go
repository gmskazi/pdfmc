package cmd

import (
	"os"
	"os/exec"
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

func TestMergeCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(tempDir string)
		args           func(dir string) []string
		fileOutput     string
		expectError    bool
		expectedOutput string
		checkFile      bool
	}{
		{
			name: "Merge two PDF files",
			setup: func(tempDir string) {
				err := createValidPDF(tempDir + "/file1.pdf")
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(tempDir + "/file2.pdf")
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}
			},
			args: func(_ string) []string {
				return []string{"merge", "file1.pdf", "file2.pdf"}
			},
			fileOutput:     "merged_output.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
		{
			name: "Merge two PDF files with custom name",
			setup: func(tempDir string) {
				err := createValidPDF(tempDir + "/file1.pdf")
				if err != nil {
					t.Fatalf("failed to create file1.pdf: %v", err)
				}
				err = createValidPDF(tempDir + "/file2.pdf")
				if err != nil {
					t.Fatalf("failed to create file2.pdf: %v", err)
				}
			},
			args: func(_ string) []string {
				return []string{"merge", "file1.pdf", "file2.pdf", "-n", "testCustomName"}
			},
			fileOutput:     "testCustomName.pdf",
			expectError:    false,
			expectedOutput: "PDF files merged successfully to:",
			checkFile:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.setup != nil {
				tt.setup(tempDir)
			}

			err := os.Chdir(tempDir)
			if err != nil {
				t.Fatalf("failed to change directory: %v", tempDir)
			}

			cmdArgs := tt.args(tempDir)
			cmd := exec.Command("pdfmc", cmdArgs...)
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				assert.Error(t, err, "Expected an error but command ran successfuly.")
			} else {
				assert.NoError(t, err, "Expected command to run successfuly but it failed.")
			}

			assert.Contains(t, string(output), tt.expectedOutput, "Unexpected output from command.")

			if tt.checkFile {
				_, err := os.Stat(tt.fileOutput)
				if err != nil {
					assert.NoError(t, err, "Expected merged PDF file to be created but it wasn't there.")
					_ = os.Remove(tt.fileOutput)
				}
			}
		})
	}
}
