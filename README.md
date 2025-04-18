# pdfmc

[![Go Report Card](https://goreportcard.com/badge/github.com/gmskazi/pdfmc)](https://goreportcard.com/report/github.com/gmskazi/pdfmc)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Workflow](https://github.com/gmskazi/pdfmc/actions/workflows/ci.yml/badge.svg)](https://github.com/gmskazi/pdfmc/actions)
[![GitHub release](https://img.shields.io/github/v/release/gmskazi/pdfmc)](https://github.com/gmskazi/pdfmc/releases/latest)

## PDF CLI tool

Pdfmc stands for PDF Merge Crypt.

A simple PDF tool to merge and encrypt files, I've created this tool to learn more about [golang](https://go.dev/),
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
> If you're using ZSH, you may need to add the $GOPATH to your ~/.zshrc.

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

> Note: **For Mac users**, you may need to allow pdfmc to run on your mac the first time you run it.
Navigate to System Settings > Privacy & Security > scroll down to the bottom and allow pdfmc to run.

## Run

### Merge PDFs

![pdfmc merge](public/merge.gif)
Navigate to the directory where your PDFs live and run:

```bash
pdfmc merge
```

Or you have the option to add a directory and it will search that directory for pdf files.

```bash
pdfmc merge ~/Downloads
```

Or you can add pdf files that you would like to be merged, this will skip the UI.

```bash
pdfmc merge file1.pdf file2.pdf file3.pdf
```

You will receive a file "merged_output.pdf", this file will be located in your current working directory and will have
all the PDFs combined into one files.

#### Flags

---

- Custom name for the merged PDF file (default "merged_output")

> '--name' or '-n' flag.

```bash
pdfmc merge -n testname
```

> Output file: testname.pdf

- Reorder the PDFs through the UI.

> '--order' or '-o' flag.

```bash
pdfmc merge -o
```

- Encrypt the PDF through the UI.

> '--encrypt' or '-e' flag.

```bash
pdfmc merge -e
```

- Set the password so its non-interactive.

> '--password' or '-p' flag.

```bash
pdfmc merge -p veryStr0ngPa33w0rd!
```

> Note: you can't use the --password and --encrypt flags together, you will need to use one or the other.

#### Merge example interactive mode

> This will merge, order and encrypt the files interactively through the UI.

```bash
pdfmc merge -eo
```

Output file

```bash
merged_output.pdf
```

#### Merge example Non-interactive mode

> This will merge, encrypt and set a custom filename non-interactively.

```bash
pdfmc merge file1.pdf file2.pdf -n testname -p veryStr0ngPa33w0rd!
```

Output file

```bash
testname.pdf
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

Or you can add pdf files that you would like to be encrypted, this will skip the UI for selecting the PDF files.

```bash
pdfmc merge file1.pdf file2.pdf file3.pdf
```

You have the ability to choose which PDFs you would like to encrypt (including multiple files) and set a password.

#### flags

---

- Password to encrypt the PDF files.

> '--password' or '-p' flag.

The below option will use the UI for selecting the pdf files to encrypt.

```bash
pdfmc encrypt -p veryStr0ngPa33w0rd!
```

#### Encrypt example interactive mode

> Encrypt and set a password interactively through the UI.

```bash
pdfmc encrypt
```

#### Encrypt example non-interactive mode

> Encrypt and set a password non-interactively.

```bash
pdfmc encrypt file1.pdf file2.pdf -p veryStr0ngPa33w0rd!
```

---

### Decrypt PDFs

![pdfmc decrypt](public/decrypt.gif)

Navigate to the directory where your PDFs live and run:

```bash
pdfmc decrypt
```

Or you have the option to add a directory and it will search that directory for pdf files.

```bash
pdfmc decrypt directory
```

Or you can add pdf files that you would like to be decrypted, this will skip the UI for selecting the PDF files.

```bash
pdfmc merge file1.pdf file2.pdf file3.pdf
```

You have the ability to choose which PDFs you would like to encrypt (including multiple files) and set a password.

#### flags

---

- Password to decrypt the PDF files.

> '--password' or '-p' flag.

The below option will use the UI for selecting the pdf files to decrypt.

```bash
pdfmc decrypt -p veryStr0ngPa33w0rd!
```

#### Decrypt example interactive mode

> Decrypt the files interactively through the UI.

```bash
pdfmc decrypt
```

#### Decrypt example non-interactive mode

> Decrypt and set a password non-interactively.

```bash
pdfmc Decrypt file1.pdf file2.pdf -p veryStr0ngPa33w0rd!
```

---

## Completions

![completions](public/completions.gif)

> Note: auto-completions will only work on MacOS or Linux.

Custom and normal completions that have been configured:

- Autocompletion for the subcommands by hitting tab.
- After 'pdfmc merge/encrypt' only pdf files and folders will be displayed.
- If you choose a pdf file as your first choice your second options shouldn't include your first choice and will only
display pdf files.
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
compinit
```

- Reload zsh

```bash
source ~/.zshrc
```

## Contributing

See [contributing](CONTRIBUTING.md).
