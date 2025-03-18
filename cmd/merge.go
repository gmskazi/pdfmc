/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gmskazi/pdfmc/cmd/pdf"
	"github.com/gmskazi/pdfmc/cmd/ui/multiReorder"
	"github.com/gmskazi/pdfmc/cmd/ui/multiSelect"
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
	Use:   "merge",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileUtils := utils.NewFileUtils()

		dir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		pdfs := fileUtils.GetPdfFilesFromDir(dir)

		for {
			p := tea.NewProgram(multiSelect.MultiSelectModel(pdfs, dir, "merge"))
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

			reorder, err := cmd.Flags().GetBool("order")
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			// reordering of the pdfs
			if reorder {
				r := tea.NewProgram(multiReorder.MultiReorderModel(selectedPdfs, "merge"))
				result, err = r.Run()
				if err != nil {
					fmt.Println(errorStyle.Render(err.Error()))
					return
				}

				reorderModel := result.(multiReorder.Tmodel)
				if reorderModel.Quit {
					os.Exit(0)
				}

				selectedPdfs = reorderModel.GetOrderedPdfs()
			}

			pdfWithFullPath := fileUtils.AddFullPathToPdfs(dir, selectedPdfs)

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
				return
			}

			pdfProcessor := pdf.NewPDFProcessor(fileUtils)

			name = pdfProcessor.PdfExtension(name)

			err = pdfProcessor.MergePdfs(pdfWithFullPath, name)
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
	mergeCmd.Flags().BoolP("order", "o", false, "Reorder the PDF files before merging.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
