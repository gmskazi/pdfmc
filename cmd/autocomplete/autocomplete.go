package autocomplete

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type FilterOptions struct {
	BaseDir   string
	HomeDir   string
	Filter    string
	Args      []string
	UsedFiles map[string]bool
}

func GetSuggestions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	homeDir, baseDir, filter, err := resolveBaseDir(toComplete)
	if err != nil {
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

	opts := FilterOptions{
		BaseDir:   baseDir,
		HomeDir:   homeDir,
		Filter:    filter,
		Args:      args,
		UsedFiles: usedFiles,
	}

	suggestions := filterPDFsAndDirs(files, opts)

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

func filterPDFsAndDirs(files []os.DirEntry, opts FilterOptions) []string {
	var suggestions []string
	for _, file := range files {
		name := file.Name()
		lower := strings.ToLower(name)
		fullPath := filepath.Join(opts.BaseDir, name)

		if opts.HomeDir != "" && strings.HasPrefix(fullPath, opts.HomeDir) {
			fullPath = "~/" + strings.TrimPrefix(fullPath, opts.HomeDir+"/")
		}

		if len(opts.Args) == 0 {
			if file.IsDir() || strings.HasSuffix(lower, ".pdf") {
				if strings.HasPrefix(lower, strings.ToLower(opts.Filter)) {
					suggestions = append(suggestions, fullPath)
				}
			}
		} else {
			if !file.IsDir() && strings.HasSuffix(lower, ".pdf") && strings.HasPrefix(lower, strings.ToLower(opts.Filter)) {
				if !opts.UsedFiles[fullPath] {
					suggestions = append(suggestions, fullPath)
				}
			}
		}
	}
	return suggestions
}
