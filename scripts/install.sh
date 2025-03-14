#!/bin/sh
# Ensure the script is running as root or with sudo
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root or with sudo."
    exec sudo "$0" "$@"
fi

set -e

REPO="gmskazi/pdfmc"
BIN_NAME="pdfmc"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="$(mktemp -d 2>/dev/null || echo '/tmp/pdfmc_install_$$')"

# Make the temp directory
mkdir -p "$TMP_DIR"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$OS" != "darwin" ]; then
    OS="linux"
fi

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
TAG_NUMBER=$(echo "$LATEST_TAG" | cut -c2-)

# Construct download URL
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/${BIN_NAME}_${TAG_NUMBER}_${OS}_${ARCH}.tar.gz"
ZIP_FILE="${BIN_NAME}_${TAG_NUMBER}_${OS}_${ARCH}.tar.gz"

# Download and install
echo "Downloading $BIN_NAME $LATEST_TAG for $OS/$ARCH..."
curl -L "$URL" -o "$TMP_DIR/$ZIP_FILE"

echo "Extracting $ZIP_FILE..."
tar -xzf "$TMP_DIR/$ZIP_FILE" -C "$TMP_DIR"

echo "Installing $ZIP_FILE..."
mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

# Clean up
rm -rf "$TMP_DIR"

if command -v "$BIN_NAME" >/dev/null 2>&1; then
    echo "$BIN_NAME installed successfully!"
else
    echo "There was an issue installing $BIN_NAME"
fi
