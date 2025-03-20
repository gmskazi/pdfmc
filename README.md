# pdfmc

[![Go Report Card](https://goreportcard.com/badge/github.com/gmskazi/pdfmc)](https://goreportcard.com/report/github.com/gmskazi/pdfmc)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Workflow](https://github.com/gmskazi/pdfmc/actions/workflows/ci.yml/badge.svg)](https://github.com/gmskazi/pdfmc/actions)
[![GitHub release](https://img.shields.io/github/v/release/gmskazi/pdfmc)](https://github.com/gmskazi/pdfmc/releases/latest)

## PDF CLI tool

Pdfmc stands for PDF Merge Crypt.

A simple PDF tool to merge and encrypt files, I've created this tool to learn more about golang,
[cobra](https://github.com/spf13/cobra), [bubbletea](https://github.com/charmbracelet/bubbletea) and
[lipgloss](https://github.com/charmbracelet/lipgloss), but after creating it, I've started using it and added it to my
toolkit. Hoping to add more functionality to it when I have the time.

## Install

### If you have golang installed

---

```bash
go install github.com/gmskazi/pdfmc@latest
```

This installs a go binary that will automatically bind to your $GOPATH.
> If you're using ZSH, you'll need to add it manually to ~/.zshrc.

```bash
GOPATH=$HOME/go PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

Reload your zshrc config

```bash
source ~/.zshrc
```

---

### Install using a script

Run the below command in the terminal of your choice to install pdfmc.

```sh
curl -L https://raw.githubusercontent.com/gmskazi/pdfmc/main/scripts/install.sh | sh
```

---

### Uninstall using a script

Run the below command in the terminal of your choice to uninstall pdfmc.

```sh
curl -L https://raw.githubusercontent.com/gmskazi/pdfmc/main/scripts/uninstall.sh | sh
```

---

## Run

### Merge PDFs

![pdfmc merge](public/merge.gif)
Navigate to the directory where your PDFs live and run:

```bash
pdfmc merge
```

Or you have the option to add a directory and it will search that directory for pdf files.

```bash
pdfmc merge directory
```

Or you can add pdf files that you would like to be merged, this will skip the UI.

```bash
pdfmc merge file1.pdf file2.pdf file3.pdf
```

You will receive a file "merged_output.pdf" that has all the PDFs combined into one files.

#### Flags

---

- Custom name for the merged PDF file (default "merged_output")

'--name' or '-n' flag.

```bash
pdfmc merge -n testname
```

- Reorder the PDFs through the UI.

'--order' or '-o' flag.

```bash
pdfmc merge -o
```

---

### Encrypt PDFs

![pdfmc encrypt](public/encrypt.gif)

Navigate to the directory where your PDFs live and run:

```bash
pdfmc encrypt
```

Or you have the option to add a directory and it will search that directory for pdf files.

```bash
pdfmc encrypt directory
```

Or you can add pdf files that you would like to be merged, this will skip the UI for selecting the PDF files.

```bash
pdfmc merge file1.pdf file2.pdf file3.pdf
```

You have the ability to choose which PDFs you would like to encrypt (including multiple files) and set a password,
the encrypted files will have "encrypt-" at the beginning of the file saved in the same directory.

#### flags

---

- Password to encrypt the PDF files.

'--password' or '-p' flag.

```bash
pdfmc encrypt -p veryStr0ngPa33w0rd!
```

Or

```bash
pdfmc encrypt file1.pdf file2.pdf -p veryStr0ngPa33w0rd!
```

## Completions

<!-- TODO: Add a gif to demonstrate autocomplete -->
![completions](public/completions.gif)

Custom and normal completions that have been configured:

- Autocompletion for the subcommands by hitting tab.
- After 'pdfmc merge/encrypt' only pdf files and folders will be displayed.
- If you choose a pdf file as your first choice your second options shouldn't include your first choice.
- If you add a '-' + tab all the flags will be displayed.

- To add Completions for zsh:

```bash
mkdir -p ~/.zsh/completions
pdfmc completion zsh > ~/.zsh/completions/_pdfmc
```

- Update your zshrc file

Add the below to your '.zshrc' file

```bash
# Custom completions
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit
compinit -i
```

- Reload zsh

```bash
source ~/.zshrc
```

## Contributing

See [contributing](CONTRIBUTING.md).
