/*
Copyright © 2025 Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiReorder"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	infoStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#5dd2fc")).Bold(true)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FC895F")).Bold(true)
	errorStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#ba0b0b")).Bold(true)
	name          string
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [files... or folder]",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			selectedPdfs []string
			quit         bool
		)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		reorder, err := cmd.Flags().GetBool("order")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		encrypt, err := cmd.Flags().GetBool("encrypt")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		pword, err := cmd.Flags().GetString("password")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		if encrypt && pword != "" {
			cmd.PrintErrln(errorStyle.Render("Please provide either the --password flag or use the --encrypt flag for interactive encryption."))
			return
		}
		fileUtils := utils.NewFileUtils(args)

		// check if any files/folders are provided
		pdfs, dir, interactive, err := fileUtils.CheckProvidedArgs()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		if interactive {
			for {
				selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, "merge")
				if err != nil || quit {
					cmd.PrintErrln(errorStyle.Render(err.Error()))
					return
				}

				if len(selectedPdfs) <= 1 {
					continue
				}
				break
			}
		}

		if !interactive {
			selectedPdfs = pdfs
		}

		// reordering of the pdfs
		if reorder {
			selectedPdfs, quit, err = multiReorder.MultiReorderInteractive(selectedPdfs, "merge")
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}
		}

		pdfWithFullPath := fileUtils.AddFullPathToPdfs(dir, selectedPdfs)

		pdfProcessor := pdf.NewPDFProcessor(fileUtils, "merge")

		name, err = pdfProcessor.MergePdfs(pdfWithFullPath, name)
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		saveDir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		// if the encrypt flag is set, ask for password interactively
		if encrypt {
			pword, quit, err = textInputs.TextinputInteractive()
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			fmt.Println()
		}

		// encrypt pdf file if flag is set
		if pword != "" {
			nonEncryptedFile := name
			fmt.Println(name)
			name, err = pdfProcessor.EncryptPdf(nonEncryptedFile, saveDir, pword)
			if err != nil {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}
			if err := os.Remove(nonEncryptedFile); err != nil {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}
		}

		complete := fmt.Sprintf("PDF files merged successfully to: %s/%s", saveDir, name)
		cmd.Println(infoStyle.Render(complete))
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&name, "name", "n", "merged_output", "Custom name for the merged PDF files")
	mergeCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF file.")
	mergeCmd.Flags().BoolP("order", "o", false, "Reorder the PDF files before merging.")
	mergeCmd.Flags().BoolP("encrypt", "e", false, "Encrypt the PDF file interatively.")

	// autocomplete for files flag
	mergeCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var suggestions []string
		var homeDir string
		dir := "."

		if strings.HasPrefix(toComplete, "/") {
			dir = "/"
			toComplete = strings.TrimPrefix(toComplete, "/")
		} else if strings.HasPrefix(toComplete, "~/") {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return nil, cobra.ShellCompDirectiveError
			}
			dir = homeDir
			toComplete = strings.TrimPrefix(toComplete, "~/")
		}

		// Determine the base directory and filter term
		searchDir := filepath.Join(dir, toComplete)
		if fileInfo, err := os.Stat(searchDir); err == nil && fileInfo.IsDir() {
			dir = searchDir
			toComplete = "" // Reset filtering because we are inside a valid folder
		} else {
			dir = filepath.Dir(searchDir) // Use the parent directory
			toComplete = filepath.Base(searchDir)
		}

		// Track already selected PDFs
		usedFiles := make(map[string]bool)
		for _, arg := range args {
			usedFiles[strings.TrimSpace(arg)] = true
		}

		// Read directory contents (1-level deep, no recursion)
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		for _, file := range files {
			name := file.Name()
			lowerName := strings.ToLower(name)
			fullPath := filepath.Join(dir, name)

			// Convert home paths back to `~/`
			if homeDir != "" && strings.HasPrefix(fullPath, homeDir) {
				fullPath = "~/" + strings.TrimPrefix(fullPath, homeDir+"/")
			}

			// If no args, suggest directories and PDFs
			if len(args) == 0 {
				if (file.IsDir() || strings.HasSuffix(lowerName, ".pdf")) && strings.HasPrefix(lowerName, strings.ToLower(toComplete)) {
					suggestions = append(suggestions, fullPath)
				}
			} else {
				// If args exist, only suggest PDFs that haven't been selected yet
				if !file.IsDir() && strings.HasSuffix(lowerName, ".pdf") && strings.HasPrefix(lowerName, strings.ToLower(toComplete)) {
					if !usedFiles[fullPath] {
						suggestions = append(suggestions, fullPath)
					}
				}
			}
		}

		if len(suggestions) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		// Prevent spaces from being inserted after completing a directory
		return suggestions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
