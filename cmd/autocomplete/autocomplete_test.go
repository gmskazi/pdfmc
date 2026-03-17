package autocomplete

import (
	"os"
	"testing"

	"github.com/gmskazi/pdfmc/cmd/testhelpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
			testhelpers.CreateTestFiles(t, tempDir, []string{"one.pdf", "two.pdf", "three.txt"})
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
