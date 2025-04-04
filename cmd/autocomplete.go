package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func GetSuggestions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	homeDir, baseDir, filter, err := resolveBaseDir(toComplete)
	if err != nil {
		cmd.PrintErrln(errorStyle.Render(err.Error()))
		return nil, cobra.ShellCompDirectiveError
	}

	files, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	usedFiles := make(map[string]bool)
	for _, arg := range args {
		usedFiles[strings.TrimSpace(arg)] = true
	}

	suggestions := filterPDFsAndDirs(files, baseDir, homeDir, filter, args, usedFiles)

	if len(suggestions) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return suggestions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
}

func resolveBaseDir(toComplete string) (homeDir, baseDir, filter string, err error) {
	dir := "."

	if strings.HasPrefix(toComplete, "/") {
		dir = "/"
		toComplete = strings.TrimPrefix(toComplete, "/")
	} else if strings.HasPrefix(toComplete, "~/") {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return "", "", "", err
		}
		dir = homeDir
		toComplete = strings.TrimPrefix(toComplete, "~/")
	}

	searchDir := filepath.Join(dir, toComplete)
	info, err := os.Stat(searchDir)
	if err == nil && info.IsDir() {
		return homeDir, searchDir, "", nil
	}

	return homeDir, filepath.Dir(searchDir), filepath.Base(searchDir), nil
}

func filterPDFsAndDirs(files []os.DirEntry, baseDir, homeDir, filter string, args []string, usedFiles map[string]bool) []string {
	var suggestions []string
	for _, file := range files {
		name := file.Name()
		lower := strings.ToLower(name)
		fullPath := filepath.Join(baseDir, name)

		if homeDir != "" && strings.HasPrefix(fullPath, homeDir) {
			fullPath = "~/" + strings.TrimPrefix(fullPath, homeDir+"/")
		}

		if len(args) == 0 {
			if file.IsDir() || strings.HasSuffix(lower, ".pdf") {
				if strings.HasPrefix(lower, strings.ToLower(filter)) {
					suggestions = append(suggestions, fullPath)
				}
			}
		} else {
			if !file.IsDir() && strings.HasSuffix(lower, ".pdf") && strings.HasPrefix(lower, strings.ToLower(filter)) {
				if !usedFiles[fullPath] {
					suggestions = append(suggestions, fullPath)
				}
			}
		}
	}
	return suggestions
}
