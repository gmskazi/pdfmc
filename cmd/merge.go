/*
Copyright © 2025 NAME HERE Aito Nakajima
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

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
	Use:   "merge [files...]",
	Short: "Merge PDFs together.",
	Long:  `This is a tool to merge PDFs together.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fileUtils := utils.NewFileUtils()

		dir, err := fileUtils.GetCurrentWorkingDir()
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		pdfs, err := cmd.Flags().GetStringSlice("files")
		if err != nil {
			fmt.Println(errorStyle.Render(err.Error()))
			return
		}

		if len(pdfs) == 0 {
			pdfs = fileUtils.GetPdfFilesFromDir(dir)
		} else {
			// Validate if the files exists when using a flag
			for _, pdf := range pdfs {
				if _, err := os.Stat(pdf); os.IsNotExist(err) {
					fmt.Println(errorStyle.Render(err.Error()))
					os.Exit(1)
				}
				if info, err := os.Stat(pdf); err == nil && info.IsDir() {
					fmt.Println(errorStyle.Render(err.Error()))
					os.Exit(1)
				}
			}
		}

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
	mergeCmd.Flags().StringSliceP("files", "f", []string{}, "File/files to merge together comma seperated.(file.pdf,file2.pdf,etc)")
	mergeCmd.Flags().BoolP("order", "o", false, "Reorder the PDF files before merging.")

	// autocomplete for files flag
	err := mergeCmd.RegisterFlagCompletionFunc("files", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var suggestions []string

		dir := "."
		parts := strings.Split(toComplete, ",")
		lastPart := strings.TrimSpace(parts[len(parts)-1])
		usedFiles := make(map[string]bool)
		for _, part := range parts[:len(parts)-1] {
			usedFiles[strings.TrimSpace(part)] = true
		}

		fmt.Fprintf(os.Stderr, "toComplete: '%s', lastPart: '%s', usedFiles: '%v'\n", toComplete, lastPart, usedFiles)

		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		for _, file := range files {
			name := file.Name()
			lowerName := strings.ToLower(name)
			if !file.IsDir() && strings.HasSuffix(lowerName, ".pdf") && strings.HasPrefix(lowerName, strings.ToLower(lastPart)) && !usedFiles[name] {
				if len(parts) > 1 {
					prefix := strings.Join(parts[:len(parts)-1], ",") + ","
					suggestions = append(suggestions, prefix+name)
				} else {
					suggestions = append(suggestions, name)
				}
			}
		}
		fmt.Fprintf(os.Stderr, "Suggestions: %v\n", suggestions)
		if len(suggestions) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		} else {
			return suggestions, cobra.ShellCompDirectiveNoSpace
		}
	})
	if err != nil {
		fmt.Println(errorStyle.Render(err.Error()))
		os.Exit(1)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
