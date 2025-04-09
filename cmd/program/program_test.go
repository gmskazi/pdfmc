package program

import (
	"testing"

	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewProgram(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
		logo string
	}
	tests := []struct {
		name string
		args args
		want *Program
	}{
		{
			name: "empty flags",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
				logo: "merge",
			},
			want: &Program{
				cmd:   &cobra.Command{},
				args:  []string{},
				logo:  "merge",
				name:  "",
				pword: "",
				MergeFlags: MergeFlags{
					reorder: false,
					encrypt: false,
				},
			},
		},
		{
			name: "all flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("name", "", "")
					cmd.Flags().Bool("order", false, "")
					cmd.Flags().Bool("encrypt", false, "")
					cmd.Flags().String("password", "", "")
					// Set the values
					cmd.SetArgs([]string{
						"--name=testName",
						"--order=true",
						"--encrypt=true",
						"--password=testPassword",
					})
					err := cmd.Execute()
					assert.NoError(t, err, "error parseing the flags")
					return cmd
				}(),
				args: []string{"file1.pdf", "file2.pdf"},
				logo: "merge",
			},
			want: &Program{
				cmd:   &cobra.Command{},
				args:  []string{"file1.pdf", "file2.pdf"},
				logo:  "merge",
				name:  "testName",
				pword: "testPassword",
				MergeFlags: MergeFlags{
					reorder: true,
					encrypt: true,
				},
			},
		},
		{
			name: "partial flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("name", "", "")
					cmd.Flags().Bool("order", false, "")
					cmd.Flags().Bool("encrypt", false, "")
					cmd.Flags().String("password", "", "")
					// Set the values
					cmd.SetArgs([]string{
						"--name=testName",
						"--order=true",
					})
					err := cmd.Execute()
					assert.NoError(t, err, "error parseing the flags")
					return cmd
				}(),
				args: []string{"file1.pdf", "file2.pdf"},
				logo: "merge",
			},
			want: &Program{
				cmd:   &cobra.Command{},
				args:  []string{"file1.pdf", "file2.pdf"},
				logo:  "merge",
				name:  "testName",
				pword: "",
				MergeFlags: MergeFlags{
					reorder: true,
					encrypt: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProgram(tt.args.cmd, tt.args.args, tt.args.logo)

			assert.Equal(t, tt.want.args, got.args, "args should match")
			assert.Equal(t, tt.want.logo, got.logo, "logo should match")
			assert.Equal(t, tt.want.MergeFlags, got.MergeFlags, "MergeFlags should match")
			assert.Equal(t, tt.want.pword, got.pword, "password should match")
		})
	}
}

func Test_getFlagValue(t *testing.T) {
	// setup test flags
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flagSet.String("stringFlag", "default", "test string flag")
	err := flagSet.Set("stringFlag", "testValue")
	assert.NoError(t, err)

	type args struct {
		flag *pflag.Flag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "flag values",
			args: args{flag: flagSet.Lookup("stringFlag")},
			want: "testValue",
		},
		{
			name: "nil flag",
			args: args{flag: nil},
			want: "",
		},
		{
			name: "unset string flag",
			args: args{flag: func() *pflag.Flag {
				fs := pflag.NewFlagSet("temp", pflag.ContinueOnError)
				fs.String("unsetString", "default", "")
				return fs.Lookup("unsetString")
			}()},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFlagValue(tt.args.flag)
			assert.Equal(t, tt.want, got, "Flags should match")
		})
	}
}

type textInputFunc func() (string, bool, error)

var TextinputInteractive textInputFunc = textInputs.TextinputInteractive

func Test_getPassword(t *testing.T) {
	tests := []struct {
		name          string
		program       *Program
		inputFunc     textInputFunc
		expectedError bool
		expectedPword string
	}{
		{
			name: "password already set",
			program: &Program{
				pword: "existingPassword",
			},
			inputFunc:     nil,
			expectedError: false,
			expectedPword: "existingPassword",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputFunc != nil {
				TextinputInteractive = tt.inputFunc
			} else {
				TextinputInteractive = textInputs.TextinputInteractive
			}

			err := tt.program.getPassword()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPword, tt.program.pword)
			}
		})
	}
}
