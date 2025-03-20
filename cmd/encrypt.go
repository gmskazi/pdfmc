/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"
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
		var selectedPdfs []string
		var interactive bool
		var pword string

		fileUtils := utils.NewFileUtils()

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
			dir := args[0]
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

		dir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
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

		pword, err = cmd.Flags().GetString("password")
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
			err := api.EncryptFile(pdf, encryptedPdfName, conf)
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			fmt.Println(selectedStyle.Render("Encrypted files: ", encryptedPdfName))
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("password", "p", "", "Password to encrypt the PDF files.")

	// autocomplete for files flag
	encryptCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var suggestions []string
		dir := "."

		// if no args, suggest directories and pdf files
		if len(args) == 0 {
			files, err := os.ReadDir(dir)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			for _, file := range files {
				name := file.Name()
				lowerName := strings.ToLower(name)
				if (file.IsDir() || strings.HasSuffix(lowerName, ".pdf")) && strings.HasPrefix(lowerName, strings.ToLower(toComplete)) {
					suggestions = append(suggestions, name)
				}
			}
		} else {
			// If args exists, assume pdf files and filter out used ones
			usedFiles := make(map[string]bool)
			for _, arg := range args {
				usedFiles[strings.TrimSpace(arg)] = true
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			for _, file := range files {
				name := file.Name()
				lowerName := strings.ToLower(name)
				if !file.IsDir() && strings.HasSuffix(lowerName, ".pdf") && strings.HasPrefix(lowerName, strings.ToLower(toComplete)) && !usedFiles[name] {
					suggestions = append(suggestions, name)
				}
			}
		}
		if len(suggestions) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return suggestions, cobra.ShellCompDirectiveDefault
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
