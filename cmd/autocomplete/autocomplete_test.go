package autocomplete

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
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

func createTestFiles(t *testing.T, tempDir string) {
	// Create test files
	files := []string{"one.pdf", "two.pdf", "three.txt"}
	for _, f := range files {
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

func TestGetSuggestions(t *testing.T) {
	type args struct {
		cmd        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 cobra.ShellCompDirective
	}{
		{
			name: "Return pdfs and directories from current directory",
			args: args{
				cmd:        &cobra.Command{},
				args:       []string{},
				toComplete: "",
			},
			want: []string{
				"one.pdf",
				"subdir",
				"two.pdf",
			},
			want1: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault,
		},
		{
			name: "Complete using 'o'",
			args: args{
				cmd:        &cobra.Command{},
				args:       []string{},
				toComplete: "o",
			},
			want: []string{
				"one.pdf",
			},
			want1: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault,
		},
		{
			name: "Returns empty when no match is found",
			args: args{
				cmd:        &cobra.Command{},
				args:       []string{},
				toComplete: "notExisting",
			},
			want:  nil,
			want1: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault,
		},
		{
			name: "Directory completion",
			args: args{
				cmd:        &cobra.Command{},
				args:       []string{},
				toComplete: "sub",
			},
			want: []string{
				"subdir",
			},
			want1: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault,
		},
		{
			name: "Handles existing files in args",
			args: args{
				cmd:        &cobra.Command{},
				args:       []string{"one.pdf"},
				toComplete: "",
			},
			want: []string{
				"two.pdf",
			},
			want1: cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			createTestFiles(t, tempDir)
			if err := os.Chdir(tempDir); err != nil {
				assert.NoError(t, err, "failed to change directory: ", tempDir)
			}
			got, got1 := GetSuggestions(tt.args.cmd, tt.args.args, tt.args.toComplete)
			assert.EqualValues(t, tt.want, got, "GetSuggestions() got = %v, want %v", got, tt.want)
			if tt.want == nil {
				assert.EqualValues(t, 4, got1, "GetSuggestions() got1 = %v, want %v", got1, 4)
			} else {
				assert.EqualValues(t, tt.want1, got1, "GetSuggestions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
