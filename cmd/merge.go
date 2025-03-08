/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gmskazi/pdfmergecrypt/cmd/ui/multiSelect"
	"github.com/gmskazi/pdfmergecrypt/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	infoStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#5dd2fc")).Bold(true)
	errorStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#ba0b0b")).Bold(true)
	name       string
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := utils.GetCurrentWorkingDir()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		pdfs, err := utils.GetPdfFilesFromDir(dir)
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		for {
			p := tea.NewProgram(multiSelect.MultiSelectModel(pdfs, dir))
			result, err := p.Run()
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			model := result.(multiSelect.Tmodel)
			if model.Quit {
				os.Exit(0)
			}
			selectedPdfs := model.GetSelectedPDFs()

			if len(selectedPdfs) <= 1 {
				continue
			}

			pdfWithFullPath := utils.AddFullPathToPdfs(dir, pdfs)

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			name = utils.PdfExtension(name)

			err = utils.MergePdfs(pdfWithFullPath, name)
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			complete := fmt.Sprintf("PDF files merged successfully to: %s/%s", dir, name)
			fmt.Println(infoStyle.Render(complete))

			break
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&name, "name", "n", "merged_output", "Custom name for the merged PDF files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
