/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"

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
	Use:   "encrypt",
	Short: "Encrypt PDF files.",
	Long:  `This is a tool for encrypting PDF files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileUtils := utils.NewFileUtils()

		dir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		pdfs := fileUtils.GetPdfFilesFromDir(dir)

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

		selectedPdfs := mSelectModel.GetSelectedPDFs()

		if len(selectedPdfs) == 0 {
			fmt.Println(infoStyle.Render("No PDFs were selected. Exiting."))
			return
		}

		p := tea.NewProgram(textInputs.TextinputModel())
		result, err = p.Run()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		tmodel := result.(textInputs.Tmodel)
		if tmodel.Quit {
			os.Exit(0)
		}

		pword := tmodel.GetPassword()

		conf := model.NewAESConfiguration(pword, pword, 256)

		fmt.Println()

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
