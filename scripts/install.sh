#!/bin/sh
set -e

REPO="gmskazi/pdfmc"
BIN_NAME="pdfmc"
INSTALL_DIR="/usr/local/bin"

# Ensure the script is running as root or with sudo
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root or with sudo."
    exec sudo "$0" "$@"
fi

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
elif [ "$ARCH" = "amd64" ]; then
    ARCH="amd64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Fetch latest release tag from GitHub API
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | cut -d '"' -f 4)

# Construct download URL
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/${BIN_NAME}-${OS}-${ARCH}"

# Download and install
echo "Downloading $BIN_NAME $LATEST_TAG for $OS/$ARCH..."
curl -L "$URL" -o "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo "$BIN_NAME installed successfully!"
