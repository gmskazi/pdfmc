/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gmskazi/pdfmergecrypt/cmd/ui/multiSelect"
	textInputs "github.com/gmskazi/pdfmergecrypt/cmd/ui/textinput"
	"github.com/gmskazi/pdfmergecrypt/cmd/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		for _, pdf := range selectedPdfs {
			encryptedPdfName := "encrypted-" + pdf
			err := api.EncryptFile(pdf, encryptedPdfName, conf)
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			fmt.Println(focusedStyle.Render("Encrypted files: ", encryptedPdfName))
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
