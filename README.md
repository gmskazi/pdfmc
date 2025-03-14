# pdfmc

## PDF CLI tool

Pdfmc stands for PDF Merge Crypt.

A simple PDF tool to merge and encrypt files, I'm creating this tool
to learn more about golang, cobra, bubbletea and lipgloss.

## Install

### If you have golang installed

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

### Install using a script

Run the below command in the terminal of your choice to install pdfmc.

```sh
curl -L https://raw.githubusercontent.com/gmskazi/pdfmc/main/scripts/install.sh | sh
```

### Uninstall using a script

Run the below command in the terminal of your choice to uninstall pdfmc.

```sh
curl -L https://raw.githubusercontent.com/gmskazi/pdfmc/main/scripts/uninstall.sh | sh
```

## Run

### Merge PDFs

![pdfmc merge](public/merge.gif)
Navigate to the directory where your PDFs live and run:

```bash
pdfmc merge
```

You will receive a file "merged_output.pdf" that has all of the PDFs
combined into one files.

#### Change the output file name

You have the ability to change the name of the output file by using the '--name'
or '-n' flag.

```bash
pdfmc merge --name testname
```

### Encrypt PDFs

![pdfmc encrypt](public/encrypt.gif)

```bash
pdfmc encrypt
```

You have the ability to choose which PDFs you would like to encrypt
(including multiple files) and set a password, the encrypted files will have
"encrypt-" at the beginning of the file saved in the same directory.

## Contributing

See [contributing](CONTRIBUTING.md).
