/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmc/cmd/ui/textinputs"
	"github.com/gmskazi/pdfmc/cmd/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [files... or folder]",
	Short: "Encrypt PDF files.",
	Long:  `This is a tool for encrypting PDF files.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var pdfs []string
		var dir string
		var selectedPdfs []string
		var interactive bool
		var pword string

		fileUtils := utils.NewFileUtils(args)

		// check if any files/folders are provided
		if len(args) == 0 {
			interactive = true
			dir, err := fileUtils.GetCurrentWorkingDir()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}
			pdfs = fileUtils.GetPdfFilesFromDir(dir)

		} else if len(args) == 1 && fileUtils.IsDirectory(args[0]) {
			interactive = true
			dir = args[0]
			pdfs = fileUtils.GetPdfFilesFromDir(dir)

		} else {
			interactive = false
			pdfs = args
			for _, pdf := range pdfs {
				if _, err := os.Stat(pdf); os.IsNotExist(err) {
					fmt.Println(errorStyle.Render(fmt.Sprintf("File '%s' does not exists", pdf)))
					os.Exit(1)
				}
				if info, err := os.Stat(pdf); err == nil && info.IsDir() {
					fmt.Println(errorStyle.Render(fmt.Sprintf("'%s' is a directory, not a file", pdf)))
					os.Exit(1)
				}
			}
		}

		if interactive {
			mSelect := tea.NewProgram(multiSelect.MultiSelectModel(pdfs, dir, "encrypt"))
			result, err := mSelect.Run()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			mSelectModel := result.(multiSelect.Tmodel)
			if mSelectModel.Quit {
				os.Exit(0)
			}

			selectedPdfs = mSelectModel.GetSelectedPDFs()

			if len(selectedPdfs) == 0 {
				fmt.Println(infoStyle.Render("No PDFs were selected. Exiting."))
				return
			}
		}

		pword, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		if pword == "" {
			p := tea.NewProgram(textInputs.TextinputModel())
			result, err := p.Run()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			tmodel := result.(textInputs.Tmodel)
			if tmodel.Quit {
				os.Exit(0)
			}

			pword = tmodel.GetPassword()

			fmt.Println()
		}

		conf := model.NewAESConfiguration(pword, pword, 256)

		if !interactive {
			selectedPdfs = pdfs
		}

		for _, pdf := range selectedPdfs {
			encryptedPdfName := "encrypted-" + pdf
			if dir == "" {
				err := api.EncryptFile(pdf, encryptedPdfName, conf)
				if err != nil {
					fmt.Println(errorStyle.Render(err.Error()))
					return
				}
			} else {
				err := api.EncryptFile(dir+"/"+pdf, encryptedPdfName, conf)
				if err != nil {
					fmt.Println(errorStyle.Render(err.Error()))
					return
				}
			}

			saveDir, err := fileUtils.GetCurrentWorkingDir()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			complete := fmt.Sprintf("PDF files encrypted successfully to: %s/%s", saveDir, encryptedPdfName)
			fmt.Println(selectedStyle.Render(complete))
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
