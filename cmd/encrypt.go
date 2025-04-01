/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [files... or folder]",
	Short: "Encrypt PDF files.",
	Long:  `This is a tool for encrypting PDF files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			selectedPdfs []string
			pword        string
			quit         bool
		)

		pword, err := cmd.Flags().GetString("password")
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		fileUtils := utils.NewFileUtils(args)
		pdfProcessor := pdf.NewPDFProcessor(fileUtils)

		// check if any files/folders are provided
		pdfs, dir, interactive, err := fileUtils.CheckProvidedArgs()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		if interactive {
			selectedPdfs, quit, err = multiSelect.MultiSelectInteractive(pdfs, dir, "encrypt")
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			if len(selectedPdfs) == 0 {
				cmd.PrintErrln(infoStyle.Render("No PDFs were selected. Exiting."))
				return
			}
		}

		if pword == "" {
			pword, quit, err = textInputs.TextinputInteractive()
			if err != nil || quit {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			fmt.Println()
		}

		if !interactive {
			selectedPdfs = pdfs
		}

		saveDir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			cmd.PrintErrln(errorStyle.Render(err.Error()))
			return
		}

		for _, pdf := range selectedPdfs {
			encryptedPdf, err := pdfProcessor.EncryptPdf(pdf, dir, pword)
			if err != nil {
				cmd.PrintErrln(errorStyle.Render(err.Error()))
				return
			}

			complete := fmt.Sprintf("PDF files encrypted successfully to: %s/%s", saveDir, encryptedPdf)
			cmd.Println(selectedStyle.Render(complete))
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF files.")

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
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
